import { NotebookPanel } from '@jupyterlab/notebook';
import { executeCode } from './notebook';
import MarkdownIt from 'markdown-it';

export interface Turn {
  type: 'llm' | 'user';
  msg: string;
}

export interface LlmVar {
  id: string;
  name: string;
  type: string;
  input: string;
}

export interface ApplicationState {
  isPythonNotebook: boolean;
  llmVars: LlmVar[];
  currentLlmVar: LlmVar | null;
  history: Turn[];
  initializing: boolean;
  thinking: boolean;
}

export function newApplicationState(
  isPythonNotebook: boolean
): ApplicationState {
  return {
    isPythonNotebook: isPythonNotebook,
    llmVars: [],
    currentLlmVar: null,
    history: [],
    initializing: true,
    thinking: false
  };
}

export class State {
  private _applicationState: ApplicationState = newApplicationState(false);
  private _applicationStateChangedCallbacks: ((
    state: ApplicationState
  ) => void)[] = [];
  private _currentNotebookPanel: NotebookPanel | null = null;
  private _isTicking: boolean = false;
  private _updateStatus: 'none' | 'sending' | 'discard' = 'none';

  private static _instance: State | null = null;

  static getInstance(): State {
    if (this._instance === null) {
      this._instance = new State();
    }
    return this._instance;
  }

  private constructor() {
    setInterval(() => this._tick(), 1000);
  }

  private async _tick() {
    if (this._isTicking) {
      return;
    }

    if (this._updateStatus === 'sending') {
      return;
    }

    if (this._updateStatus === 'discard') {
      this._updateStatus = 'none';
    }

    if (this._applicationState.isPythonNotebook && this._currentNotebookPanel) {
      this._isTicking = true;
      try {
        let newState = { ...this._applicationState };

        await this._queryLlmVars(newState);
        this._setCurrentLlmVar(newState);
        await this._queryHistory(newState);

        newState.initializing = false;

        if ((this._updateStatus as any) === 'sending') {
          return;
        }

        if ((this._updateStatus as any) === 'discard') {
          this._updateStatus = 'none';
          return;
        }

        this.applicationState = newState;
      } catch (e) {
        console.log('error', e);
      } finally {
        this._isTicking = false;
      }
    }
  }

  private async exec(output: { [key: string]: string[] }) {
    let payload: { [key: string]: string } = {};
    for (let key in output) {
      let fun = output[key][0];
      let vars = output[key].slice(1);

      vars = vars.map(
        v =>
          "'" +
          v.replace(/'/g, "\\'").replace(/\r/g, '').replace(/\n/g, '\\n') +
          "'"
      );
      let funCall =
        '__langforge_jupyterlab__helpers__instance__.' +
        fun +
        '(' +
        vars.join(',') +
        ')';
      payload[key] = funCall;
    }
    let pythonResult = await executeCode(
      this._currentNotebookPanel!,
      '',
      payload
    );
    let result: { [key: string]: any } = {};
    for (let key in pythonResult) {
      result[key] = this._resultToJson(pythonResult[key]);
    }
    return result;
  }

  private async _queryLlmVars(newState: ApplicationState) {
    try {
      newState.llmVars = (
        await this.exec({ llmVars: ['get_llm_vars'] })
      ).llmVars;
    } catch (e) {
      const error = e as Error;
      error.message = 'queryLlmVars: ' + error.message;
      throw error;
    }
  }

  private _setCurrentLlmVar(newState: ApplicationState) {
    if (
      newState.currentLlmVar !== null &&
      newState.llmVars.filter(
        llmVar => llmVar.name === newState.currentLlmVar!.name
      ).length === 0
    ) {
      newState.currentLlmVar = null;
    }

    if (!newState.currentLlmVar && newState.llmVars.length > 0) {
      newState.currentLlmVar = newState.llmVars[0];
    }
  }

  private async _queryHistory(newState: ApplicationState) {
    try {
      if (!newState.currentLlmVar) {
        return;
      }
      newState.history = (
        await this.exec({
          history: [
            'get_history',
            newState.currentLlmVar.name,
            newState.currentLlmVar.input
          ]
        })
      ).history;
      newState.history = newState.history.map(h => markdown(h));
    } catch (e) {
      const error = e as Error;
      error.message = 'queryHistory: ' + error.message;
      throw error;
    }
  }

  private _resultToJson(result: any): any {
    if (result && result.status === 'ok') {
      let dataTextPlain = result.data['text/plain'];

      // remove leading and trailing single quotes if present
      if (dataTextPlain.startsWith("'") && dataTextPlain.endsWith("'")) {
        dataTextPlain = dataTextPlain.substring(1, dataTextPlain.length - 1);
        // replace escaped single quotes
        dataTextPlain = dataTextPlain.replace(/\\'/g, "'");
        // escape double quotes
        dataTextPlain = dataTextPlain.replace(/\\\\"/g, '\\"');
      }
      result = JSON.parse(dataTextPlain);

      return result;
    }
    return null;
  }

  onApplicationStateChanged(callback: (state: ApplicationState) => void) {
    this._applicationStateChangedCallbacks.push(callback);
  }

  get applicationState() {
    return this._applicationState;
  }

  set applicationState(state: ApplicationState) {
    this._applicationState = state;
    this._applicationStateChangedCallbacks.forEach(callback => callback(state));
  }

  get currentNotebookPanel() {
    return this._currentNotebookPanel;
  }

  set currentNotebookPanel(currentNotebookPanel: NotebookPanel | null) {
    this._currentNotebookPanel = currentNotebookPanel;
  }

  async sendMessage(msg: string) {
    if (
      !this._applicationState.isPythonNotebook ||
      !this._currentNotebookPanel ||
      !this._applicationState.currentLlmVar
    ) {
      return;
    }
    this._updateStatus = 'sending';

    try {
      const newState = { ...this._applicationState };
      newState.history.push(markdown({ msg: msg, type: 'user' }));
      newState.thinking = true;

      this.applicationState = newState;

      newState.history = (
        await this.exec({
          history: [
            'send_message',
            newState.currentLlmVar!.name,
            newState.currentLlmVar!.input,
            msg
          ]
        })
      ).history;

      newState.history = newState.history.map(h => markdown(h));
      newState.thinking = false;

      this.applicationState = newState;
    } finally {
      this._updateStatus = 'discard';
    }
  }

  selectLlmVar(id: string) {
    const newState = { ...this._applicationState };
    newState.currentLlmVar = newState.llmVars.filter(
      llmVar => llmVar.id === id
    )[0];
    this.applicationState = newState;
  }
}

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
});

function markdown(turn: Turn): Turn {
  let text = turn.msg;
  try {
    // escape double quotes
    text = text.replace(/"/g, '\\"');
    text = JSON.parse(`"${text}"`);
  } catch (err) {
    console.error('Invalid string format:', err);
    console.error(text);
    return turn;
  }

  if (turn.type === 'user') {
    return { type: turn.type, msg: text.replace(/\n/g, '<br/>') };
  } else {
    text = text.trim();
    text = md.render(text);
    text = text
      .replace(/<pre>/g, '<pre class="jp-LangForge">')
      .replace(/<code>/g, '<code class="jp-LangForge">');
    text = text.trim();

    // remove leading and trailing <p> tags
    if (text.startsWith('<p>') && text.endsWith('</p>')) {
      text = text.substring(3, text.length - 4);
    }

    return { type: turn.type, msg: text.replace(/\n/g, '<br/>') };
  }
}

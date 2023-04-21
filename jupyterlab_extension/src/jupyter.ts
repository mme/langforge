import { NotebookPanel } from '@jupyterlab/notebook';
import { executeCode } from './notebook';
import { newJupyterState, JupyterState } from './interfaces';
import { markdown } from './markdown-utils';

export class Jupyter {
  private _jupyterState: JupyterState = newJupyterState(false);
  private _jupyterStateChangedCallbacks: ((state: JupyterState) => void)[] = [];
  private _currentNotebookPanel: NotebookPanel | null = null;
  private _isTicking: boolean = false;
  private _updateStatus: 'none' | 'sending' | 'discard' = 'none';

  private static _instance: Jupyter | null = null;

  static getInstance(): Jupyter {
    if (this._instance === null) {
      this._instance = new Jupyter();
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

    if (this._jupyterState.isPythonNotebook && this._currentNotebookPanel) {
      this._isTicking = true;
      try {
        let newState = { ...this._jupyterState };

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

        this.state = newState;
      } catch (e) {
        console.log('error', e);
      } finally {
        this._isTicking = false;
      }
    }
  }

  private async _exec(output: { [key: string]: string[] }) {
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

  private async _queryLlmVars(newState: JupyterState) {
    try {
      newState.llmVars = (
        await this._exec({ llmVars: ['get_llm_vars'] })
      ).llmVars;
    } catch (e) {
      const error = e as Error;
      error.message = 'queryLlmVars: ' + error.message;
      throw error;
    }
  }

  private _setCurrentLlmVar(newState: JupyterState) {
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

  private async _queryHistory(newState: JupyterState) {
    try {
      if (!newState.currentLlmVar) {
        return;
      }
      newState.history = (
        await this._exec({
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

  onStateChanged(callback: (state: JupyterState) => void) {
    this._jupyterStateChangedCallbacks.push(callback);
  }

  get state() {
    return this._jupyterState;
  }

  set state(state: JupyterState) {
    this._jupyterState = state;
    this._jupyterStateChangedCallbacks.forEach(callback => callback(state));
  }

  get currentNotebookPanel() {
    return this._currentNotebookPanel;
  }

  set currentNotebookPanel(currentNotebookPanel: NotebookPanel | null) {
    this._currentNotebookPanel = currentNotebookPanel;
  }

  async sendMessage(msg: string) {
    if (
      !this._jupyterState.isPythonNotebook ||
      !this._currentNotebookPanel ||
      !this._jupyterState.currentLlmVar
    ) {
      return;
    }
    this._updateStatus = 'sending';

    try {
      const newState = { ...this._jupyterState };
      newState.history.push(markdown({ text: msg, type: 'input' }));
      newState.thinking = true;

      this.state = newState;

      newState.history = (
        await this._exec({
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

      this.state = newState;
    } finally {
      this._updateStatus = 'discard';
    }
  }

  selectLlmVar(id: string) {
    const newState = { ...this._jupyterState };
    newState.currentLlmVar = newState.llmVars.filter(
      llmVar => llmVar.id === id
    )[0];
    this.state = newState;
  }
}

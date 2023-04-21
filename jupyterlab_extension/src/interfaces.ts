export interface Turn {
  type: 'input' | 'output';
  text: string;
}

export interface LlmVar {
  id: string;
  name: string;
  type: string;
  input: string;
}

export interface JupyterState {
  isPythonNotebook: boolean;
  llmVars: LlmVar[];
  currentLlmVar: LlmVar | null;
  history: Turn[];
  initializing: boolean;
  thinking: boolean;
}

export function newJupyterState(
  isPythonNotebook: boolean,
  initializing: boolean = true
): JupyterState {
  return {
    isPythonNotebook: isPythonNotebook,
    llmVars: [],
    currentLlmVar: null,
    history: [],
    initializing: initializing,
    thinking: false
  };
}

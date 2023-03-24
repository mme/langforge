import { NotebookPanel } from '@jupyterlab/notebook';
import { KernelMessage } from '@jupyterlab/services';
import { JSONObject } from '@lumino/coreutils';

export async function executeCode(
  notebookPanel: NotebookPanel,
  code: string,
  userExpressions: JSONObject = {}
): Promise<any> {
  if (!notebookPanel) {
    throw new Error('Notebook is null or undefined.');
  }

  await notebookPanel.sessionContext.ready;

  if (!notebookPanel.sessionContext.session) {
    throw new Error('Notebook session is null or undefined.');
  }

  if (!notebookPanel.sessionContext.session.kernel) {
    throw new Error('Notebook session kernel is null or undefined.');
  }

  const kernel = notebookPanel.sessionContext.session.kernel;
  const message: KernelMessage.IShellMessage = await kernel.requestExecute({
    allow_stdin: true,
    code: code,
    silent: true,
    stop_on_error: false,
    store_history: false,
    user_expressions: userExpressions
  }).done;

  const content: any = message.content;
  if (content.status !== 'ok') {
    throw new Error(content.ename + ': ' + content.evalue);
  }

  return content.user_expressions;
}

import {
  ILabShell,
  ILayoutRestorer,
  JupyterFrontEnd,
  JupyterFrontEndPlugin
} from '@jupyterlab/application';
import { LangForgeSidebarWidget } from './sidebar';
import { NotebookPanel } from '@jupyterlab/notebook';
import { Widget } from '@lumino/widgets';
// import { executeCode } from './notebook';
import { Jupyter } from './jupyter';
import { newJupyterState } from './interfaces';
import { ILauncher } from '@jupyterlab/launcher';
import { addLaunchers } from './launcher-utils';
import { monkeyPatchLauncher } from './monkey-patch';

/**
 * Initialization data for the langforge extension.
 */
const plugin: JupyterFrontEndPlugin<void> = {
  id: 'langforge:plugin',
  autoStart: true,
  requires: [ILayoutRestorer, ILabShell, ILauncher],
  activate: (
    app: JupyterFrontEnd,
    restorer: ILayoutRestorer,
    labShell: ILabShell,
    launcher: ILauncher
  ) => {
    const { commands } = app;
    monkeyPatchLauncher();
    addLaunchers(commands, launcher);

    const state = Jupyter.getInstance();
    // Instantiate the widget
    const sidebarWidget = new LangForgeSidebarWidget(state);

    // Add the widget to the right sidebar
    app.shell.add(sidebarWidget, 'right', { rank: 1 });

    // Track the widget state
    restorer.add(sidebarWidget, sidebarWidget.id);

    app.restored.then(() => {
      currentChanged(labShell.currentWidget);

      labShell.currentChanged.connect(async (sender, args) => {
        currentChanged(args.newValue);
      });
    });
  }
};

async function currentChanged(newValue: Widget | null) {
  const state = Jupyter.getInstance();
  if (newValue instanceof NotebookPanel) {
    await newValue.sessionContext.ready;
    let kernelName = newValue.sessionContext.session?.kernel?.name;
    let isPythonNotebook =
      kernelName !== null &&
      kernelName !== undefined &&
      kernelName.startsWith('python');
    state.state = newJupyterState(isPythonNotebook, false);
    if (isPythonNotebook) {
      state.currentNotebookPanel = newValue;
    } else {
      state.currentNotebookPanel = null;
    }
  } else {
    state.state = newJupyterState(false, false);
    state.currentNotebookPanel = null;
  }
}

export default plugin;

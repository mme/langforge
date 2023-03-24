import {
  ILabShell,
  ILayoutRestorer,
  JupyterFrontEnd,
  JupyterFrontEndPlugin
} from '@jupyterlab/application';
import { LangForgeSidebarWidget } from './sidebar';
import { INotebookTracker, NotebookPanel } from '@jupyterlab/notebook';
import { Widget } from '@lumino/widgets';
// import { executeCode } from './notebook';
import { newApplicationState, State } from './state';

/**
 * Initialization data for the langforge extension.
 */
const plugin: JupyterFrontEndPlugin<void> = {
  id: 'langforge:plugin',
  autoStart: true,
  requires: [ILayoutRestorer, INotebookTracker, ILabShell],
  activate: (
    app: JupyterFrontEnd,
    restorer: ILayoutRestorer,
    notebookTracker: INotebookTracker,
    labShell: ILabShell
  ) => {
    const state = State.getInstance();
    // Instantiate the widget
    const sidebarWidget = new LangForgeSidebarWidget(state);

    // Add the widget to the right sidebar
    app.shell.add(sidebarWidget, 'right', { rank: 1 });

    // Track the widget state
    restorer.add(sidebarWidget, sidebarWidget.id);

    app.restored.then(() => {
      labShell.currentChanged.connect(async (sender, args) => {
        const newValue = args.newValue as Widget;
        if (newValue instanceof NotebookPanel) {
          await newValue.sessionContext.ready;
          let kernelName = newValue.sessionContext.session?.kernel?.name;
          let isPythonNotebook =
            kernelName !== null &&
            kernelName !== undefined &&
            kernelName.startsWith('python');
          state.applicationState = newApplicationState(isPythonNotebook);
          if (isPythonNotebook) {
            state.currentNotebookPanel = newValue;
          } else {
            state.currentNotebookPanel = null;
          }
        } else {
          state.applicationState = newApplicationState(false);
          state.currentNotebookPanel = null;
        }
      });
      // notebookTracker.currentChanged.connect(async (tracker, panel) => {
      //   state;
      //   console.log('currentChanged', tracker, panel);
      //   if (panel) {
      //     let notebookName: string = panel.context.path;
      //     let result = await executeCode(panel, 'pass', {
      //       globals:
      //         '{k: v for k, v in globals().items() if not k.startswith("_")}'
      //     });
      //     console.log(result.globals.data['text/plain']);
      //     console.log('session Context', panel.sessionContext.isReady);
      //     let isPython = notebookName.endsWith('.ipynb');
      //     if (isPython) {
      //       console.log('Python notebook');
      //     } else {
      //       console.log('Not a Python notebook');
      //     }
      //   } else {
      //     console.log('No active notebook');
      //   }
      // });
    });
  }
};

export default plugin;

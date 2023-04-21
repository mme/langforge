import { Panel } from '@lumino/widgets';
import { LabIcon } from '@jupyterlab/ui-components';
import * as React from 'react';
// import { useState } from 'react';
import { Widget } from '@lumino/widgets';
import { ReactWidget } from '@jupyterlab/apputils';
import { Chat } from './chat';
import { Jupyter } from './jupyter';
import { JupyterState, newJupyterState } from './interfaces';
import { langForgeSvgStr } from './svg';

const chatIcon = new LabIcon({
  name: 'custom:chat',
  svgstr: langForgeSvgStr
});

const handleSendMessage = (message: string) => {
  Jupyter.getInstance().sendMessage(message);
};

const handleSelectLlmVar = (id: string) => {
  Jupyter.getInstance().selectLlmVar(id);
};

function Root(props: { state: Jupyter }): JSX.Element {
  const [applicationState, setApplicationState] = React.useState<JupyterState>(
    newJupyterState(false)
  );

  React.useEffect(() => {
    props.state.onStateChanged((state: JupyterState) => {
      setApplicationState(state);
    });
  }, []);

  let showChat = false;
  if (
    !applicationState.initializing &&
    applicationState.isPythonNotebook &&
    applicationState.llmVars.length > 0
  ) {
    showChat = true;
  }

  return (
    <div className="jp-LangForge-root">
      {!showChat && (
        <div
          style={{
            flex: 1,
            justifyContent: 'center',
            alignItems: 'center',
            display: 'flex',
            textAlign: 'center',
            padding: '10px'
          }}
        >
          {applicationState.initializing && <>Initializing...</>}
          {!applicationState.initializing &&
            !applicationState.isPythonNotebook && (
              <> No Python notebook selected</>
            )}
          {!applicationState.initializing &&
            applicationState.isPythonNotebook &&
            applicationState.llmVars.length === 0 && (
              <> No LLM variables found in the notebook</>
            )}
        </div>
      )}

      {applicationState.isPythonNotebook &&
        applicationState.llmVars.length > 0 && (
          <div className="jp-LangForge-content">
            <Chat
              currentLlmVar={applicationState.currentLlmVar!}
              history={applicationState.history}
              llmVars={applicationState.llmVars}
              onSendMessage={handleSendMessage}
              onSelectLlmVar={handleSelectLlmVar}
              thinking={applicationState.thinking}
            />
          </div>
        )}
    </div>
  );
}

export class LangForgeSidebarWidget extends Panel {
  // private dockPanel: DockPanel;

  constructor(state: Jupyter) {
    super();
    this.id = 'langforge-sidebar';
    this.title.icon = chatIcon;
    this.title.caption = 'LangForge';
    this.title.closable = true;
    this.addClass('jp-LangForge-Sidebar');

    const myWidget: Widget = ReactWidget.create(<Root state={state} />);
    myWidget.addClass('jp-LangForge-react');
    this.addWidget(myWidget);
  }
}

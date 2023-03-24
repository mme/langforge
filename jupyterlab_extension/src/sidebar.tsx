import { Panel } from '@lumino/widgets';
import { LabIcon } from '@jupyterlab/ui-components';
import * as React from 'react';
// import { useState } from 'react';
import { Widget } from '@lumino/widgets';
import { ReactWidget } from '@jupyterlab/apputils';
import { Chat } from './chat';
import { ApplicationState, newApplicationState, State } from './state';

const iconSvgStr = `
<svg width="20px" height="20px" viewBox="0 0 20 20" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <g id="Chat-interface" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g id="Group-2" transform="translate(2.000000, 2.000000)">
            <rect id="Rectangle" stroke="#404040" stroke-width="1.5" x="-0.75" y="-0.75" width="17.5" height="17.5" rx="2"></rect>
            <path d="M5.43884698,8.47077047 L5.43884698,9.2971444 C5.85578305,10.8109734 6.758037,11.5678879 8.14560884,11.5678879 C9.53318068,11.5678879 10.4354346,10.8109734 10.8523707,9.2971444 L10.8523707,8.47077047 L5.43884698,8.47077047 Z" id="Path-4" fill="#484543"></path>
            <path d="M1.34428879,7.06532866 C2.00942888,6.21035381 2.77783764,5.78286638 3.64951509,5.78286638 C4.52119253,5.78286638 5.19921875,6.21035381 5.68359375,7.06532866" id="Path-6" stroke="#0C0C0C"></path>
            <path d="M10.6418373,7.06532866 C11.2139458,6.21035381 11.909393,5.78286638 12.7281789,5.78286638 C13.5469648,5.78286638 14.1911818,6.21035381 14.6608297,7.06532866" id="Path-7" stroke="#1A1C1C"></path>
        </g>
    </g>
</svg>
`;

const chatIcon = new LabIcon({
  name: 'custom:chat',
  svgstr: iconSvgStr
});

const handleSendMessage = (message: string) => {
  State.getInstance().sendMessage(message);
};

const handleSelectLlmVar = (id: string) => {
  State.getInstance().selectLlmVar(id);
};

function Root(props: { state: State }): JSX.Element {
  const [applicationState, setApplicationState] =
    React.useState<ApplicationState>(newApplicationState(false));

  React.useEffect(() => {
    props.state.onApplicationStateChanged((state: ApplicationState) => {
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

  constructor(state: State) {
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

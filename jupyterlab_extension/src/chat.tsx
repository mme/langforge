import { useState, useRef, useEffect } from 'react';
import * as React from 'react';
import TextareaAutosize from 'react-textarea-autosize';
import { LlmVar, Turn } from './state';

interface ChatProps {
  llmVars: LlmVar[];
  currentLlmVar: LlmVar;
  history: Turn[];
  thinking: boolean;
  onSendMessage: (message: string) => void;
  onSelectLlmVar: (botId: string) => void;
}

export function Chat(props: ChatProps): JSX.Element {
  const [newMessage, setNewMessage] = useState('');

  const formRef = useRef<HTMLFormElement>(null);
  const chatMessagesRef = useRef<HTMLDivElement>(null);
  const chatInputTextAreaRef = useRef<HTMLTextAreaElement>(null);

  const handleInputChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    setNewMessage(event.target.value);
  };

  const handleFormSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    submit();
  };

  const submit = () => {
    if (props.thinking) {
      return;
    }
    props.onSendMessage(newMessage);
    setNewMessage('');
    chatInputTextAreaRef.current!.focus();
  };

  const handleSelectLlmVar = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedId = event.target.value;
    props.onSelectLlmVar(selectedId);
  };

  useEffect(() => {
    if (chatMessagesRef.current) {
      chatMessagesRef.current.scrollTop = chatMessagesRef.current.scrollHeight;
    }
  }, [
    props.currentLlmVar.id + '_' + props.currentLlmVar.input,
    props.history.length
  ]);

  const handleKeyDown: React.KeyboardEventHandler<
    HTMLTextAreaElement
  > = event => {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      submit();
    }
  };

  return (
    <div className="jp-LangForge-chat">
      <div className="jp-LangForge-chat-header">
        <div className="jp-LangForge-llm-var-select">
          <select onChange={handleSelectLlmVar} value={props.currentLlmVar.id}>
            {props.llmVars.map((llmVar, index) => {
              return (
                <option key={index} value={llmVar.id}>
                  {llmVar.name} - {llmVar.input}
                </option>
              );
            })}
          </select>
        </div>
      </div>
      <div className="jp-LangForge-chat-messages" ref={chatMessagesRef}>
        {props.history.map((turn, index) => {
          let rowClass = 'jp-LangForge-chat-message-row';
          if (turn.type === 'user') {
            rowClass += ' jp-LangForge-chat-message-row-user';
          } else if (turn.type === 'llm') {
            rowClass += ' jp-LangForge-chat-message-row-llm';
          }

          let avatarDiv = <div />;

          if (turn.type === 'user') {
            avatarDiv = (
              <div className="jp-LangForge-avatar-user">{svgAvatarUser}</div>
            );
          } else if (turn.type === 'llm') {
            avatarDiv = (
              <div className="jp-LangForge-avatar-bot">{svgAvatarLlm}</div>
            );
          }
          return (
            <div key={index} className={rowClass}>
              <div className="jp-LangForge-chat-message-container">
                <div className="jp-LangForge-chat-avatar">{avatarDiv}</div>
                <div
                  className="jp-LangForge-chat-message"
                  dangerouslySetInnerHTML={{ __html: turn.msg }}
                />
              </div>
            </div>
          );
        })}
        {props.thinking === true && (
          <div className="jp-LangForge-chat-message-row jp-LangForge-chat-message-row-llm">
            <div className="jp-LangForge-chat-message-container">
              <div className="jp-LangForge-chat-avatar">
                <div className="jp-LangForge-avatar-bot">{svgAvatarLlm}</div>
              </div>
              <div
                className="jp-LangForge-chat-message"
                style={{ marginTop: '2px' }}
                dangerouslySetInnerHTML={{ __html: svgThinking }}
              />
            </div>
          </div>
        )}
      </div>
      <form
        className="jp-LangForge-chat-input-form"
        onSubmit={handleFormSubmit}
        ref={formRef}
      >
        <div className="jp-LangForge-chat-input-container">
          <TextareaAutosize
            placeholder="Type a message..."
            value={newMessage}
            onChange={handleInputChange}
            onKeyDown={handleKeyDown}
            className="jp-LangForge-chat-input"
            maxRows={5}
            ref={chatInputTextAreaRef}
          />
          <button
            className="jp-LangForge-chat-send-button"
            disabled={props.thinking}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
              />
            </svg>
          </button>
        </div>
      </form>
    </div>
  );
}

const svgAvatarUser = (
  <svg
    width="25px"
    height="25px"
    viewBox="0 0 25 25"
    version="1.1"
    xmlns="http://www.w3.org/2000/svg"
  >
    <g
      id="Chat-interface"
      stroke="none"
      stroke-width="1"
      fill="none"
      fill-rule="evenodd"
    >
      <g id="Group-4" transform="translate(1.000000, 1.000000)">
        <rect
          id="Rectangle"
          stroke="#3A3937"
          stroke-width="0.827586207"
          fill="#D0D5F7"
          x="-0.413793103"
          y="-0.413793103"
          width="23.8275862"
          height="23.8275862"
          rx="1.44827586"
        ></rect>
        <g
          id="U"
          transform="translate(8.162109, 6.958008)"
          fill="#121826"
          fill-rule="nonzero"
        >
          <path
            d="M1.38085938,0 L1.38085938,6.20703125 C1.38085938,6.93619792 1.51757812,7.54231771 1.79101562,8.02539062 C2.19661458,8.75455729 2.88020833,9.11914062 3.84179688,9.11914062 C4.99479167,9.11914062 5.77864583,8.7249349 6.19335938,7.93652344 C6.41666667,7.50813802 6.52832031,6.93164062 6.52832031,6.20703125 L6.52832031,0 L7.90917969,0 L7.90917969,5.63964844 C7.90917969,6.87467448 7.74283854,7.82486979 7.41015625,8.49023438 C6.79947917,9.70247396 5.64648438,10.3085938 3.95117188,10.3085938 C2.25585938,10.3085938 1.10514323,9.70247396 0.499023438,8.49023438 C0.166341146,7.82486979 0,6.87467448 0,5.63964844 L0,0 L1.38085938,0 Z"
            id="Path"
          ></path>
        </g>
      </g>
    </g>
  </svg>
);

const svgAvatarLlm = (
  <svg
    width="25px"
    height="25px"
    viewBox="0 0 25 25"
    version="1.1"
    xmlns="http://www.w3.org/2000/svg"
  >
    <g
      id="Chat-interface"
      stroke="none"
      stroke-width="1"
      fill="none"
      fill-rule="evenodd"
    >
      <g id="Group-5" transform="translate(-11.000000, -18.000000)">
        <g id="Group-2" transform="translate(12.000000, 19.000000)">
          <g
            id="Group-Copy"
            fill="#A8E2E4"
            stroke="#3A3937"
            stroke-width="0.827586207"
          >
            <rect
              id="Rectangle"
              x="-0.413793103"
              y="-0.413793103"
              width="23.8275862"
              height="23.8275862"
              rx="1.44827586"
            ></rect>
          </g>
          <path
            d="M7.81834254,12.1767326 L7.81834254,13.3646451 C8.41768813,15.5407743 9.71467818,16.6288389 11.7093127,16.6288389 C13.7039472,16.6288389 15.0009373,15.5407743 15.6002829,13.3646451 L15.6002829,12.1767326 L7.81834254,12.1767326 Z"
            id="Path-4"
            fill="#484543"
          ></path>
          <path
            d="M11.0790511,16.3019935 C11.0790511,15.2918968 11.4725047,14.7868484 12.259412,14.7868484 C13.0463194,14.7868484 13.439773,15.1985677 13.439773,16.0220063 L12.4104425,16.3019935 L11.0790511,16.3019935 Z"
            id="Path-5"
            fill="#EBA7AA"
          ></path>
          <circle
            id="Oval"
            fill="#D8D8D8"
            cx="20.5215517"
            cy="13.9784483"
            r="1.09051724"
          ></circle>
          <circle
            id="Oval-Copy-2"
            fill="#D8D8D8"
            cx="3.46982759"
            cy="13.9784483"
            r="1.09051724"
          ></circle>
          <path
            d="M1.93241514,10.15641 C2.88855401,8.9273836 3.99314161,8.31287042 5.24617794,8.31287042 C6.49921426,8.31287042 7.47387695,8.9273836 8.17016602,10.15641"
            id="Path-6"
            stroke="#0C0C0C"
            stroke-width="0.724137931"
          ></path>
          <path
            d="M15.2976411,10.15641 C16.120047,8.9273836 17.1197524,8.31287042 18.2967571,8.31287042 C19.4737619,8.31287042 20.3998238,8.9273836 21.0749428,10.15641"
            id="Path-7"
            stroke="#1A1C1C"
            stroke-width="0.724137931"
          ></path>
        </g>
      </g>
    </g>
  </svg>
);

const svgThinking =
  '<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><style>.spinner_S1WN{animation:spinner_MGfb .8s linear infinite;animation-delay:-.8s}.spinner_Km9P{animation-delay:-.65s}.spinner_JApP{animation-delay:-.5s}@keyframes spinner_MGfb{93.75%,100%{opacity:.2}}</style><circle class="spinner_S1WN" cx="4" cy="12" r="3"/><circle class="spinner_S1WN spinner_Km9P" cx="12" cy="12" r="3"/><circle class="spinner_S1WN spinner_JApP" cx="20" cy="12" r="3"/></svg>';

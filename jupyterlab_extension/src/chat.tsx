import { useState, useRef, useEffect } from 'react';
import * as React from 'react';
import TextareaAutosize from 'react-textarea-autosize';
import { LlmVar, Turn } from './interfaces';
import { svgAvatarUserStr, svgAvatarLlmStr, svgThinkingStr } from './svg';

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
          if (turn.type === 'input') {
            rowClass += ' jp-LangForge-chat-message-row-user';
          } else if (turn.type === 'output') {
            rowClass += ' jp-LangForge-chat-message-row-llm';
          }

          let avatarDiv = <div />;

          if (turn.type === 'input') {
            avatarDiv = (
              <div
                className="jp-LangForge-avatar-user"
                dangerouslySetInnerHTML={{ __html: svgAvatarUserStr }}
              />
            );
          } else if (turn.type === 'output') {
            avatarDiv = (
              <div
                className="jp-LangForge-avatar-bot"
                dangerouslySetInnerHTML={{ __html: svgAvatarLlmStr }}
              />
            );
          }
          return (
            <div key={index} className={rowClass}>
              <div className="jp-LangForge-chat-message-container">
                <div className="jp-LangForge-chat-avatar">{avatarDiv}</div>
                <div
                  className="jp-LangForge-chat-message"
                  dangerouslySetInnerHTML={{ __html: turn.text }}
                />
              </div>
            </div>
          );
        })}
        {props.thinking === true && (
          <div className="jp-LangForge-chat-message-row jp-LangForge-chat-message-row-llm">
            <div className="jp-LangForge-chat-message-container">
              <div className="jp-LangForge-chat-avatar">
                <div
                  className="jp-LangForge-avatar-bot"
                  dangerouslySetInnerHTML={{ __html: svgAvatarLlmStr }}
                />
              </div>
              <div
                className="jp-LangForge-chat-message"
                style={{ marginTop: '2px' }}
                dangerouslySetInnerHTML={{ __html: svgThinkingStr }}
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

import { ILauncher } from '@jupyterlab/launcher';
import { CommandRegistry } from '@lumino/commands';
import chatCreativeNotebook from 'chat-creative.ipynb';
import chatDeterministicNotebook from 'chat-deterministic.ipynb';
import apiAgentNotebook from 'api-agent.ipynb';
import qaTextNotebook from 'qa-txt.ipynb';
import qaPdfNotebook from 'qa-pdf.ipynb';
import codeNotebook from 'code.ipynb';
import babyAgiNotebook from 'baby-agi.ipynb';

import { Contents } from '@jupyterlab/services';
import { NotebookPanel } from '@jupyterlab/notebook';

async function openNotebookWithContent(
  commands: CommandRegistry,
  content: string,
  desiredName: string
): Promise<void> {
  const model: Contents.IModel = await commands.execute(
    'docmanager:new-untitled',
    {
      path: '',
      type: 'notebook',
      ext: '.ipynb'
    }
  );

  const doc: NotebookPanel = await commands.execute('docmanager:open', {
    path: model.path
  });

  await doc.context.ready;
  doc.context.model.fromString(content);
  doc.context.model.initialize();

  let renamed = false;
  let index = 0;
  while (!renamed) {
    try {
      await doc.context.rename(
        `${desiredName}${index === 0 ? '' : '-' + index}.ipynb`
      );
      renamed = true;
    } catch {
      index += 1;
    }
  }

  await doc.context.save();
  await doc.sessionContext.ready;
  doc.content.activeCellIndex = 1;
}

function addItem(
  commands: CommandRegistry,
  launcher: ILauncher,
  category: string,
  command: string,
  label: string,
  caption: string,
  iconClass: string,
  rank: number,
  notebook: string,
  desiredName: string
) {
  commands.addCommand(command, {
    label: label,
    caption: caption,
    iconClass: iconClass,
    execute: () => {
      openNotebookWithContent(commands, notebook, desiredName);
    }
  });
  launcher.add({
    command,
    category,
    rank
  });
}

export function addLaunchers(commands: CommandRegistry, launcher: ILauncher) {
  addItem(
    commands,
    launcher,
    'Templates',
    'chat-creative:create',
    'Creative ChatGPT',
    'Open a notebook for creative chat',
    'jp-NotebookIcon',
    1,
    chatCreativeNotebook,
    'chat-creative'
  );

  addItem(
    commands,
    launcher,
    'Templates',
    'chat-deterministic:create',
    'Deterministic ChatGPT',
    'Open a notebook for deterministic chat',
    'jp-NotebookIcon',
    2,
    chatDeterministicNotebook,
    'chat-deterministic'
  );

  addItem(
    commands,
    launcher,
    'Templates',
    'api-agent:create',
    'API Agent',
    'Open a notebook for an API agent',
    'jp-NotebookIcon',
    3,
    apiAgentNotebook,
    'api-agent'
  );

  addItem(
    commands,
    launcher,
    'Templates',
    'qa-txt:create',
    'QA Text',
    'Open a notebook for QA text',
    'jp-NotebookIcon',
    4,
    qaTextNotebook,
    'qa-txt'
  );

  addItem(
    commands,
    launcher,
    'Templates',
    'qa-pdf:create',
    'QA PDF',
    'Open a notebook for QA PDF',
    'jp-NotebookIcon',
    5,
    qaPdfNotebook,
    'qa-pdf'
  );

  addItem(
    commands,
    launcher,
    'Templates',
    'code:create',
    'QA Code',
    'Open a notebook for QA code',
    'jp-NotebookIcon',
    6,
    codeNotebook,
    'code'
  );

  addItem(
    commands,
    launcher,
    'Templates',
    'baby-agi:create',
    'Baby AGI',
    'Open a notebook for Baby AGI',
    'jp-NotebookIcon',
    7,
    babyAgiNotebook,
    'baby-agi'
  );
}

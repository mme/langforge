import MarkdownIt from 'markdown-it';
import { Turn } from './interfaces';

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
});

export function markdown(turn: Turn): Turn {
  let text = turn.text;
  try {
    // escape double quotes
    text = text.replace(/"/g, '\\"');
    text = JSON.parse(`"${text}"`);
  } catch (err) {
    console.error('Invalid string format:', err);
    console.error(text);
    return turn;
  }

  if (turn.type === 'input') {
    return { type: turn.type, text: text.replace(/\n/g, '<br/>') };
  } else {
    text = text.trim();
    text = md.render(text);
    text = text
      .replace(/<pre>/g, '<pre class="jp-LangForge">')
      .replace(/<code>/g, '<code class="jp-LangForge">');
    text = text.trim();

    // remove leading and trailing <p> tags
    if (text.startsWith('<p>') && text.endsWith('</p>')) {
      text = text.substring(3, text.length - 4);
    }

    return { type: turn.type, text: text.replace(/\n/g, '<br/>') };
  }
}

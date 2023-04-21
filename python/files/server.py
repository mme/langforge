import sys
from jupyter_notebook_parser import JupyterNotebookParser # type: ignore
from flask import Flask, jsonify, request # type: ignore
import langchain.chains.base # type: ignore
import copy
from langchain.schema import ( # type: ignore
    AIMessage,
    HumanMessage,
)
from dotenv import load_dotenv # type: ignore
import os
import argparse
from waitress import serve # type: ignore

parser = argparse.ArgumentParser(description="LangForge server script")
parser.add_argument("filename", help="File name")
parser.add_argument("--port", type=int, default=2204, help="Port number (default: 2204)")
args = parser.parse_args()

filename = args.filename
port = args.port

env_path = os.path.join(os.getcwd(), '.env')
load_dotenv(env_path)

parsed = JupyterNotebookParser(filename)

source = []

cells = parsed.get_code_cell_sources()
for cell in cells:
    source.append(cell.raw_source)

code = "\n".join(source)
code = "\n".join([line for line in code.split('\n') if not line.startswith('%')])

exec(code, globals(), locals())

app = Flask(__name__)

@app.route('/chat/<name>', methods=['POST'])
def chat(name):
    found = False
    var = None

    if name in globals():
        var = globals()[name]
        if isinstance(var, langchain.chains.base.Chain) or issubclass(type(var), langchain.chains.base.Chain):
            found = True
             
    if not found:
        return jsonify({"error": 'Variable %s not found' % name}), 404    

    data = request.get_json()
    if data is None:
        return jsonify({"error": "Invalid JSON or no JSON provided"}), 400

    if not isinstance(data, dict):
        return jsonify({"error": "JSON data should be an object"}), 400

    for k, v in data.items():
        # if v is an array
        if isinstance(v, list) and k == "memory":
            for el in v:
                if not isinstance(el, str):
                    return jsonify({"error": "Invalid input %s" % k}), 400
            continue
        elif not isinstance(v, str):
            return jsonify({"error": "Invalid input %s" % k}), 400
        if k not in var.input_keys:
            return jsonify({"error": "Invalid input %s" % k}), 400
        
    var = copy.deepcopy(var)

    if 'memory' in data:
        if hasattr(var, 'memory'):
            messages = [HumanMessage(content=el) if i % 2 == 0 else AIMessage(content=el) for i, el in enumerate(data['memory'])]
            var.memory.chat_memory.messages = messages

    args = {}
    for k, v in data.items():
        if k == 'memory':
            continue
        args[k] = v
    
    result = var(args)
    json_result = dict()
    for k, v in result.items():
        if isinstance(v, str):
            json_result[k] = v
    return jsonify(json_result)

print("Running on all addresses (0.0.0.0), port %s, filename %s" % (port, filename))
serve(app, host='0.0.0.0', port=port)

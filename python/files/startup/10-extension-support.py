class __langforge_jupyterlab__helpers__:

    _history = {}

    def __init__(self):
        self._monkey_patch()

    def _history_key(self, obj, input):
        return str(id(obj)) + '_' + input
            
    def _add_to_history(self, obj, input, msg, type):
        history_key = self._history_key(obj, input)
        if history_key not in self._history:
            self._history[history_key] = []
        self._history[history_key].append({'msg': msg, 'type': type})

    def _get_type(self, obj):
        try:
            import langchain # type: ignore
            if isinstance(obj, langchain.LLMChain) or issubclass(type(obj), langchain.LLMChain):
                return 'llm_chain'
        except ImportError as e:
            pass
        return 'unknown'
        
    def _get_inputs(self, obj):
        if self._get_type(obj) == 'llm_chain':
            return [input for input in obj.input_keys if input not in ['history', 'sys']]
        else:
            raise Exception('Unknown type %s' % type)   

    def get_llm_vars(self):
        import json

        llm_chain_vars = []
        for var_name, var_val in globals().items():
            var_type = self._get_type(var_val)
            if var_type == 'unknown':
                continue
            inputs = self._get_inputs(var_val)
            for input in inputs:
                llm_chain_vars.append({'id': str(id(var_val)), 'name': var_name, 'type': var_type, 'input': input})             

        return json.dumps(llm_chain_vars) 

    _llm_chain_original_predict = None
 
    def _monkey_patch_llm_chain(self):
        try:
            from langchain import LLMChain # type: ignore
            if self._llm_chain_original_predict is None:
                self._llm_chain_original_predict = LLMChain.predict
            this = self
            def wrapped(*args, **kwargs):
                input = list(kwargs.keys())[0]
                this._add_to_history(args[0], input, kwargs[input], 'user')            
                result = this._llm_chain_original_predict(*args, **kwargs)
                this._add_to_history(args[0], input, result, 'llm')
                return result
            LLMChain.predict = wrapped

        except ImportError as e:
            return

    def _monkey_patch(self):
        self._monkey_patch_llm_chain()

    def get_history(self, name, input):
        import json

        if name not in globals():
            raise Exception('Variable %s not found' % name)

        obj = globals()[name]
        history_key = self._history_key(obj, input)
        if history_key not in self._history:
            return []

        return json.dumps(self._history[history_key])
    
    def send_message(self, name, input, msg):
        import json
        if name not in globals():
            raise Exception('Variable %s not found' % name)

        obj = globals()[name]

        if self._get_type(obj) == 'llm_chain':            
            kwargs = {}
            kwargs[input] = msg
            obj.predict(**kwargs)
            return self.get_history(name, input)
        else:
            raise Exception('Unknown type %s' % type)

__langforge_jupyterlab__helpers__instance__ = __langforge_jupyterlab__helpers__()
__langforge_jupyterlab__helpers__instance__._monkey_patch()

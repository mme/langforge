class __langforge_jupyterlab__helpers__:

    def __init__(self):
        __langforge_jupyterlab__helpers__._monkey_patch()
        __langforge_jupyterlab__helpers__._install_handlers()

    def get_llm_vars(self):
        import json

        langchain_vars = []
        for var_name, var_val in globals().items():
            if var_name.startswith('_'):
                continue
            var_type = self._get_type(var_val)
            if var_type == 'unknown':
                continue
            inputs = self._get_inputs(var_val)
            for input in inputs:
                langchain_vars.append({'id': self._get_id(var_val), 'name': var_name, 'type': var_type, 'input': input})             

        return json.dumps(langchain_vars) 
    
    def get_history(self, name, input):
        import json

        if name not in globals():
            raise Exception('Variable %s not found' % name)

        obj = globals()[name]
        if self._get_type(obj) == 'langchain':
            return json.dumps(getattr(obj, '_langforge_history'))
        else:
            raise Exception('Unknown type %s' % type)
    
    def send_message(self, name, input, msg):
        if name not in globals():
            raise Exception('Variable %s not found' % name)

        obj = globals()[name]
        if self._get_type(obj) == 'langchain':  
            kwargs = {}
            kwargs[input] = msg
            obj(kwargs)
            return self.get_history(name, input)
        else:
            raise Exception('Unknown type %s' % type)

    _langchain_original_call = None
    _langchain_original_init = None
    LANGCHAIN_SYSTEM_INPUTS = ['history', 'sys', 'chat_history']

    @staticmethod    
    def _langchain_patch():
        try:
            import langchain # type: ignore
            from langchain.chains.base import Chain # type: ignore
            langchain.verbose = True
            if __langforge_jupyterlab__helpers__._langchain_original_call is None:
                __langforge_jupyterlab__helpers__._langchain_original_call = Chain.__call__
            if __langforge_jupyterlab__helpers__._langchain_original_init is None:
                __langforge_jupyterlab__helpers__._langchain_original_init = Chain.__init__
            Chain.__call__ = __langforge_jupyterlab__helpers__._langchain_patched_call
            Chain.__init__ = __langforge_jupyterlab__helpers__._langchain_patched_init
        except ImportError as e:
            pass

    @staticmethod
    def _langchain_patched_init(self, *args, **kwargs):
        import uuid
        __langforge_jupyterlab__helpers__._langchain_original_init(self, *args, **kwargs)
        object.__setattr__(self, '_langforge_id', str(uuid.uuid4()))
        object.__setattr__(self, '_langforge_history', [])

    @staticmethod
    def _langchain_patched_call(self, inputs, return_only_outputs = False):
        inputs = self.prep_inputs(inputs)
        self.callback_manager.on_chain_start(
            {"name": self.__class__.__name__, "_langforge_obj": self},
            inputs,
            verbose=self.verbose,
        )
        try:
            outputs = self._call(inputs)
        except (KeyboardInterrupt, Exception) as e:
            self.callback_manager.on_chain_error(e, verbose=self.verbose)
            raise e
        self.callback_manager.on_chain_end(outputs, verbose=self.verbose)
        return self.prep_outputs(inputs, outputs, return_only_outputs)
    
    def _monkey_patch():
        __langforge_jupyterlab__helpers__._langchain_patch()

    def _install_handlers():
        __langforge_jupyterlab__helpers__._install_langchain_handler()

    def _get_type(self, obj):
        try:
            import langchain.chains.base # type: ignore
            if isinstance(obj, langchain.chains.base.Chain) or issubclass(type(obj), langchain.chains.base.Chain):
                return 'langchain'
        except ImportError as e:
            pass
        return 'unknown'
    
    def _get_inputs(self, obj):
        if self._get_type(obj) == 'langchain':
            return [input for input in obj.input_keys if input not in self.LANGCHAIN_SYSTEM_INPUTS]
        else:
            raise Exception('Unknown type %s' % type)   
    
    def _get_id(self, obj):
        if self._get_type(obj) == 'langchain':
            return getattr(obj, '_langforge_id')
        else:
            raise Exception('Unknown type %s' % type)
    
    from langchain.callbacks.base import BaseCallbackHandler # type: ignore
    class LangForgeCallbackHandler(BaseCallbackHandler):
        def __init__(self, color= None):
            self.indent = 0
            self.obj = None

        def on_chain_start(self, serialized, inputs, **kwargs):            
            if self.indent == 0:
                self.obj = serialized['_langforge_obj']
                langforge_history = getattr(self.obj, '_langforge_history')
                input_key = None
                for k in inputs.keys():
                    if not k in __langforge_jupyterlab__helpers__.LANGCHAIN_SYSTEM_INPUTS:
                        input_key = k
                        break
                if input_key is None:
                    input_key = list(input.keys())[0]
                langforge_history.append({"type": "input", "text": inputs[input_key]})                
            self.indent += 1

        def on_chain_end(self, outputs, **kwargs):
            self.indent -= 1
            if self.indent == 0 and len(list(outputs.values())) > 0:
                text = list(outputs.values())[0]
                langforge_history = getattr(self.obj, '_langforge_history')
                langforge_history.append({"type": "output", "text": text})

        def on_chain_error(self, error, **kwargs):
            self.indent -= 1
            if self.indent == 0:
                text = str(error)
                langforge_history = getattr(self.obj, '_langforge_history')
                langforge_history.append({"type": "output", "text": text})
            
        def on_llm_start(self, serialized, prompts, **kwargs):
            pass

        def on_llm_end(self, response, **kwargs):
            pass

        def on_llm_new_token(self, token, **kwargs):
            pass

        def on_llm_error(self, error, **kwargs):
            pass

        def on_tool_start(self, serialized, input_str, **kwargs):
            pass

        def on_agent_action(self, action, color = None, **kwargs):
            pass
        
        def on_tool_end(self, output, color = None, observation_prefix = None, llm_prefix = None, **kwargs):
            pass

        def on_tool_error(self, error, **kwargs):
            pass
        
        def on_text(self, text, color = None, end = "", **kwargs):
            pass

        def on_agent_finish(self, finish, color = None, **kwargs):
            pass

    @staticmethod
    def _install_langchain_handler():
        from langchain.callbacks import get_callback_manager # type: ignore
        from langchain.callbacks.stdout import StdOutCallbackHandler # type: ignore
        callback = get_callback_manager()
        callback.set_handlers([__langforge_jupyterlab__helpers__.LangForgeCallbackHandler(), StdOutCallbackHandler()])

__langforge_jupyterlab__helpers__instance__ = __langforge_jupyterlab__helpers__()

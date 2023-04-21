from IPython.core.magic import register_line_magic # type: ignore
@register_line_magic
def setup(line):
    "setup integrations"

    import subprocess
    import yaml  # type: ignore
    import sys
    import os
    from IPython.display import display, HTML # type: ignore
    from IPython import paths # type: ignore


    integrations_to_install = line.split()

    ipython_dir = paths.get_ipython_dir()
    startup_dir = os.path.join(ipython_dir, "profile_default", "startup")
    yaml_file_path = os.path.join(startup_dir, "integrations.yaml")

    with open(yaml_file_path, 'r') as f:
      all_integrations = yaml.safe_load(f)

    integrations = [i for i in all_integrations if i['name'] in integrations_to_install]

    def install_integrations():
        
        # Get list of installed packages
        result = subprocess.run(['pip', 'list', '--format=freeze'], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        installed_packages = result.stdout.decode('utf-8').split('\n')
        installed_packages = [pkg.split("==")[0] for pkg in installed_packages]

        write_requirements = False

        for integration in integrations:
            packages = integration.get('packages', [])
            installed = all(package in installed_packages for package in packages)

            if installed:                
                continue
            
            print(f"Installing {integration['title']} ...")

            pre_install_commands = integration.get('preInstallCommands', [])
            for command in pre_install_commands:
                result = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                if result.returncode != 0:
                    print(result.stdout.decode('utf-8'))
                    print(result.stderr.decode('utf-8'), file=sys.stderr)
                    raise RuntimeError(f"Pre-install command '{command}' failed with return code {result.returncode}")

            for package in packages:
                if not package in installed_packages:
                    write_requirements = True
                    result = subprocess.run(['pip', 'install', package], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                    if result.returncode != 0:
                        print(result.stdout.decode('utf-8'))
                        print(result.stderr.decode('utf-8'), file=sys.stderr)
                        raise RuntimeError(f"Package installation '{package}' failed with return code {result.returncode}")

            post_install_commands = integration.get('postInstallCommands', [])
            for command in post_install_commands:
                result = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                if result.returncode != 0:
                    print(result.stdout.decode('utf-8'))
                    print(result.stderr.decode('utf-8'), file=sys.stderr)
                    raise RuntimeError(f"Post-install command '{command}' failed with return code {result.returncode}")

            print(f"{integration['title']} installed ✅")
        if write_requirements:

            with open("requirements.txt", "a") as f:
                output = subprocess.check_output(["pip", "freeze", "--local"])
                f.write(output.decode("utf-8"))

    install_integrations()

    def check_api_keys(variables):
      def show_env_variables(variables):
        import os
        import ipywidgets as widgets  # type: ignore
        from IPython.display import display  # type: ignore
        from dotenv import load_dotenv, set_key, unset_key, dotenv_values  # type: ignore

          
        def on_ok_button_click(button, variable, input_box, row):
          os.environ[variable] = input_box.value
          env_vars = dotenv_values()
          if input_box.value:
              set_key(".env", variable, input_box.value)
          elif variable in env_vars:
              unset_key(".env", variable)
          update_rows()
          row.children = [row.children[0], row.children[1], edit_buttons[variable]]

        def on_edit_button_click(button, variable, row):
          input_box = widgets.Text(placeholder=f"New value for {variable}", layout=widgets.Layout(width="100%"), continuous_update=False)
          ok_button = widgets.Button(description="OK", layout=widgets.Layout(width="50px"))
          cancel_button = widgets.Button(description="Cancel", layout=widgets.Layout(width="80px"), style={'button_color': 'white', 'font_size': '10px'})

          def submit_value(change):
              if change['name'] == 'value':
                  on_ok_button_click(button, variable, input_box, row)

          def cancel_edit():
              row.children = [row.children[0], row.children[1], edit_buttons[variable]]

          input_box.observe(submit_value, names='value')
          ok_button.on_click(lambda btn: on_ok_button_click(btn, variable, input_box, row))
          cancel_button.on_click(lambda btn: cancel_edit())

          row.children = [row.children[0], row.children[1], widgets.HBox([input_box, ok_button, cancel_button], layout=widgets.Layout(width="40%"))]
          input_box.focus()

        def update_rows():
          for variable, row in rows.items():
            row[1].value = '✅' if os.environ.get(variable) else '❌'

        # Create header row
        header = widgets.HBox([
          widgets.Label(value="API key", layout=widgets.Layout(width="40%")),
          widgets.Label(value="Is set?", layout=widgets.Layout(width="20%")),
          widgets.Label(value="", layout=widgets.Layout(width="40%"))
        ])

        rows = {}
        edit_buttons = {}

        for variable in variables:
          var_name = widgets.Label(value=variable, layout=widgets.Layout(width="40%"))
          var_name.style.font_family = "monospace"        

          status = widgets.Label(value='✅' if os.environ.get(variable) else '❌', layout=widgets.Layout(width="20%"))
          edit_button = widgets.Button(description="Edit", layout=widgets.Layout(width="40%"))
          edit_buttons[variable] = edit_button
          edit_button.on_click(lambda btn, var=variable, row=None: on_edit_button_click(btn, var, btn.row))
          row = widgets.HBox([
              var_name,
              status,
              edit_button
          ])
          rows[variable] = (row, status)
          edit_button.row = row

        update_rows()
        
        table_layout = widgets.Layout(max_width='800px')
        vbox = widgets.VBox([header] + [row[0] for row in rows.values()], layout=table_layout)
        
        # Apply CSS styles
        vbox.layout.flex_flow = 'column nowrap'
        vbox.layout.align_items = 'stretch'
        header.layout.background_color = '#888'
        header.layout.color = 'white'
        header.layout.border = '1px solid #666'
        
        for row, status in rows.values():
          row.layout.border_bottom = '1px solid #666'
          row.layout.border_left = '1px solid #666'
          row.layout.border_right = '1px solid #666'
          status.layout.width = '20%'
          status.layout.text_align = 'center'
          status.layout.padding = '0 5px'

        display(vbox)

      def check_missing_variables(variables):
        import os
        class MissingEnvironmentVariableError(BaseException):
          def __init__(self, missing_api_keys):
              self.missing_api_keys = missing_api_keys

          def _render_traceback_(self):
              missing_keys = ", ".join(self.missing_api_keys)
              message = f"Please set the missing API keys using the editor: {missing_keys}"
              display(HTML(f'<div style="background-color: red; color: white; max-width: 800px; padding: 10px;">{message}</div>'))
        
        missing_variables = [var for var in variables if not os.environ.get(var)]
        
        if missing_variables:
            raise MissingEnvironmentVariableError(missing_variables)
        
      show_env_variables(variables)
      check_missing_variables(variables)

    api_keys = set()
    for integration in integrations:
        api_keys.update(integration.get('apiKeys', []))

    check_api_keys(list(api_keys))
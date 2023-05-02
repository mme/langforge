import unittest
import subprocess

class TestLangforgeCLI(unittest.TestCase):
    def test_help_output(self):
        command = "langforge -h"
        result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, shell=True)

        # Check the return code
        self.assertEqual(result.returncode, 0, f"Expected return code 0, but got {result.returncode}")

        # Check if the output contains the expected text
        expected_text = "Usage:\n  langforge [command]"
        self.assertIn(expected_text, result.stdout, f"Expected output to contain '{expected_text}', but got '{result.stdout.strip()}'")

        # Check for no errors
        self.assertEqual(result.stderr.strip(), "", f"Expected no errors, but got '{result.stderr.strip()}'")

if __name__ == "__main__":
    unittest.main()
  
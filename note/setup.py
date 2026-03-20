#!/bin/python
try:
    import os, pathlib, sys
except ImportError as e:
    print(f"Okay, error has occured while importing packages: {e}")
    sys.exit(0)

class NoInput(Exception):
    pass

class FileStructure(Exception):
    pass

class SetUp:
    def __init__(self, file_loc):
        if not file_loc or file_loc == "":
            raise NoInput("Please enter location of file")
        self.file_loc = pathlib.Path(file_loc)
        if not self.file_loc.parent.exists():
            raise FileStructure("The path is not correct")
        if not self.file_loc.parent.is_dir():
            raise FileStructure("The path is not correct")
        if not os.access(self.file_loc.parent, os.W_OK):
            raise PermissionError("Write permission has not been given")
        if self.file_loc.exists():
            raise FileStructure("The file is already created.")
        self.BASE_DIR = pathlib.Path(__file__).resolve().parent
        self.ENV_LOC = self.BASE_DIR / '.env'

    def is_env(self):
        try:
            if self.ENV_LOC.exists():
                with open(self.ENV_LOC, "r") as env:
                    if env.read().strip():
                        print("The env has already been configured")
                        return False
            print("Creating a new env file")
            return True
        except Exception as e:
            print(f"The error is in is_env: {e}")

    def env_setup(self):
        try:
            if not self.is_env():
                return False
            with open(self.ENV_LOC, "w") as env:
                write_res = env.write(f"FILE_NAME={self.file_loc}")
            if not write_res:
                return False
            return True
        except PermissionError as e:
            print(f"There is an error in env_setup: {e}")
            return False

    def save_file(self):
        try:
            with open(self.file_loc, "w") as file:
                res = file.write("")
            if not res:
                print("Something went wrong while creating the file")
                return False
            return True
        except Exception as e:
            print(f'There is an error in save_file: {e}')

    def main(self):
        methods = [
            ("Env setup", self.env_setup),
            ("File setup", self.save_file)
        ]
        for method_name, method in methods:
            print(f"Starting {method_name}")
            error_occured = False
            if method():
                print(f"{method_name} was successfully completed")
                continue
            else:
                error_occured = True
        if not error_occured:
            print("The setup wasn't fully successfull.")
        print("The set up ran perfectly")
if __name__ == "__main__":
    file_name = input("Enter the location of the app to save your notes at: ").strip()
    try:
        setup = SetUp(file_name)
        setup.main()
    except NoInput as e:
        print(e)
    except FileStructure as e:
        print(e)
    except PermissionError as e:
        print(e)
    except Exception as e:
        print(f"There is an error during initiation: {e}")




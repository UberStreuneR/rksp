import sys
import time
import logging
import os
from watchdog.observers import Observer
from watchdog.events import LoggingEventHandler, FileSystemEventHandler
from watchdog.utils.dirsnapshot import DirectorySnapshot, DirectorySnapshotDiff

# if __name__ == "__main__":
#     logging.basicConfig(level=logging.INFO,
#                         format='%(asctime)s - %(message)s',
#                         datefmt='%Y-%m-%d %H:%M:%S')
#     path = sys.argv[1] if len(sys.argv) > 1 else '.'
#     event_handler = LoggingEventHandler()
#     dir = DirectorySnapshot()
#     observer = Observer()
#     observer.schedule(event_handler, path, recursive=True)
#     observer.start()
#     try:
#         while True:
#             time.sleep(1)
#     except KeyboardInterrupt:
#         observer.stop()
#     observer.join()

class Watcher:

    def __init__(self, directory=".", handler=FileSystemEventHandler()):
        self.observer = Observer()
        self.handler = handler
        self.directory = directory

    def run(self):
        self.observer.schedule(
            self.handler, self.directory, recursive=True)
        self.observer.start()
        print("\nWatcher Running in {}/\n".format(self.directory))
        try:
            while True:
                time.sleep(1)
        except:
            self.observer.stop()
        self.observer.join()
        print("\nWatcher Terminated\n")


class MyHandler(FileSystemEventHandler):

    files = {} # filename: ["line1", "line2"]

    def __init__(self):
        dir_files = os.listdir()
        for f in dir_files:
            self.update_file("./"+f)
            print("added file", f)
        super().__init__()

    def update_file(self, filename: str, out: bool = False):
        if not ".txt" in filename:
            return
        if filename not in self.files:
            self.files[filename] = []
        old_lines = self.files[filename]
        file = open(filename, "r")
        lines = file.readlines()
        m = min(len(old_lines), len(lines))
        self.files[filename] = lines
        if not out:
            return
        for i in range(m):
            if old_lines[i] != lines[i]:
                print("Deleted:", i, old_lines[i])
                print("Added:", i, lines[i])
        if len(old_lines) < len(lines):
            for i in range(len(old_lines), len(lines)):
                print("Added:", i, lines[i])
        else:
            for i in range(len(old_lines), len(lines)):
                print("Removed:", i, old_lines[i])

    def on_modified(self, event):
        print(event) # Your code here
        if event.is_directory:
            return
        self.update_file(event.src_path, True)
    
    def on_created(self, event):
        print(event)

if __name__=="__main__":
    w = Watcher(".", MyHandler())
    w.run()
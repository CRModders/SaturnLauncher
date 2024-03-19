import webview as wv

def main():
    wv.create_window(title='Cosmos Launcher', url='win/index.html', width=800, height=600, resizable=False)
    wv.start()

if __name__ == "__main__":
    main()

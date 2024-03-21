const backgroundSelect = document.getElementById('background-select');
const themeSelect = document.getElementById('theme-select');

window.GetSettings().then(settings => {
    backgroundSelect.value = settings.backgroundType;
    themeSelect.value = settings.theme;
});

backgroundSelect.addEventListener('change', () => {
    const newBackgroundType = backgroundSelect.value;

    window.GetSettings().then(settings => {
        settings.backgroundType = newBackgroundType;

        window.SaveSettings(settings).then(() => {
        }).catch(error => {
        });
    });
});

themeSelect.addEventListener('change', () => {
    const newTheme = themeSelect.value;

    window.GetSettings().then(settings => {
        settings.theme = newTheme;

        window.SaveSettings(settings).then(() => {
        }).catch(error => {
        });
    });
});
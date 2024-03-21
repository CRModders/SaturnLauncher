const themeSel = document.getElementById('theme-select');
const body = document.getElementsByTagName('body')[0];
const pageWrapper = document.getElementById('page-wrapper');
const sideMenu = document.getElementsByClassName('sideMenu')[0];
const exitBtnLight = document.getElementById('exitBtnLightMode')[0];
const exitBtnDark = document.getElementById('exitBtnDarkMode')[0];
const newProfileMenu = document.getElementsByClassName('newProfileScreen')[0];
const contentMenu = document.getElementsByClassName('contentMenu')[0];

// Assuming you have elements with these IDs
const emptyText = document.getElementById('emptyText');
const emptyTextGuide = document.getElementById('emptyTextGuide');
const title = document.getElementById('title');
const backgroundSelectText = document.getElementById('background-select-text');
const themeSelectText = document.getElementById('theme-select-text');

themeSel.addEventListener('change', () => {
    const newTheme = themeSel.value;

    if (newTheme === "light") {
        lightMode();
    } else {
        darkMode();
    }
});

function lightMode() {
    body.style.backgroundColor = "rgb(230, 230, 230)";
    pageWrapper.style.backgroundColor = "rgb(230, 230, 230)";
    sideMenu.style.backgroundColor = "rgba(200, 200, 200, 0.2)";
    sideMenu.style.backdropFilter = "blur(5px)";
    exitBtnLight.style.display = "block";
    exitBtnDark.style.display = "none";
    emptyText.style.color = "rgb(0, 0, 0)";
    emptyTextGuide.style.color = "rgb(50, 50, 50)";
    title.style.color = "rgb(0, 0, 0)";
    backgroundSelectText.style.color = "rgb(0, 0, 0)";
    backgroundSelectText.style.backgroundColor = "rgb(150, 150, 150)";
    themeSelectText.style.color = "rgb(0, 0, 0)";
    themeSelectText.style.backgroundColor = "rgb(150, 150, 150)";
    newProfileMenu.style.backgroundColor = "rgba(0, 0, 0, 0.2)";
    contentMenu.style.backgroundColor = "rgba(200, 200, 200, 0.3)";
}

function darkMode() {
    body.style.backgroundColor = "rgb(25, 25, 25)";
    pageWrapper.style.backgroundColor = "rgb(25, 25, 25)";
    sideMenu.style.backgroundColor = "rgba(40, 40, 40, 0.3)";
    sideMenu.style.backdropFilter = "blur(10px)";
    exitBtnLight.style.display = "none";
    exitBtnDark.style.display = "block";
    emptyText.style.color = "rgb(255, 255, 255)";
    emptyTextGuide.style.color = "rgb(200, 200, 200)";
    title.style.color = "rgb(255, 255, 255)";
    backgroundSelectText.style.color = "rgb(255, 255, 255)";
    backgroundSelectText.style.backgroundColor = "rgb(70, 70, 70)";
    themeSelectText.style.color = "rgb(255, 255, 255)";
    themeSelectText.style.backgroundColor = "rgb(70, 70, 70)";
    newProfileMenu.style.backgroundColor = "rgba(0, 0, 0, 0.2)";
    contentMenu.style.backgroundColor = "rgba(40, 40, 40, 0.3)";
}
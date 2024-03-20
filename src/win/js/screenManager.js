const profileScreen = document.getElementById('profileScreen');
const settingsScreen = document.getElementById('settingsScreen');
const newProfileScreen = document.getElementById('newProfileScreen');

profileScreen.style.display = 'block';
settingsScreen.style.display = 'none';
newProfileScreen.style.display = 'none';

function setSettingsScreen() {
    settingsScreen.style.display = 'block';
    profileScreen.style.display = 'none';
    newProfileScreen.style.display = 'none';
}

function setProfileScreen() {
    profileScreen.style.display = 'block';
    settingsScreen.style.display = 'none';
    newProfileScreen.style.display = 'none';
}

function exitNewProfileScreen() {
    newProfileScreen.style.display = 'none';
}

function newProfile() {
    newProfileScreen.style.display = 'block';
}
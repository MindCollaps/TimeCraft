document.getElementById("btnSwitchColorMode").addEventListener("change", (e) => {
    if (document.documentElement.getAttribute('data-bs-theme') == 'dark') {
        document.documentElement.setAttribute('data-bs-theme', 'light')
        storeCookie("colorMode", "light", 365)
    } else {
        document.documentElement.setAttribute('data-bs-theme', 'dark')
        storeCookie("colorMode", "dark", 365)
    }
});

loadColorMode();

function loadColorMode() {
    let colorMode = getCookie("colorMode")
    if (colorMode == "light") {
        document.documentElement.setAttribute('data-bs-theme', 'light')
        document.getElementById("btnSwitchColorMode").checked = false;
    } else {
        document.documentElement.setAttribute('data-bs-theme', 'dark')
        document.getElementById("btnSwitchColorMode").checked = true;
    }
}
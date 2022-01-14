const countdownEl = document.getElementById("countdown");
const statusBtnEl = document.getElementById("status-button");
const camerasDivEl = document.getElementById("cams");

function disableStatusButton() {
    statusBtnEl.disabled = true;
    statusBtnEl.style.backgroundColor = "rgb(178, 49, 49)";
    statusBtnEl.style.borderColor     = "rgb(94, 14, 14)";
}

function enableStatusButton() {
    statusBtnEl.disabled = false;
    statusBtnEl.style.backgroundColor = "rgb(19, 154, 19)";
    statusBtnEl.style.borderColor     = "rgb(10, 78, 10)";
}

async function get(url = "") {
    const response = await fetch(url);
    return response.json();
}

async function post(url = "", data = {}) {
    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });
    return response.json();
}

async function recover() {
    try {
        await post("/reset-timer");
    } catch {
        console.log("timer have already reset");
    } finally {
        countdown();
        disableStatusButton();
    }
}

window.onload = async () => {
    await render()
    countdown()
};

statusBtnEl.onclick = async () => await recover();

function countdown() {
    setTimeout(async function again() {
        const timeLeft = await get("/countdown");
        updateCountdown(timeLeft);
        if (!timeLeft) {
            enableStatusButton();
            return
        }

        setTimeout(again, 1000);
    });
}

async function render() {
    const cameraIDs = await get("/cameras-ids");

    for (let idx = 0; idx < cameraIDs.length; idx++) {
        const states = await get(`/cameras-info/${cameraIDs[idx]}`)
        createCamera(`CAM-${idx}`, states);
    }
}

function createCamera(name) {
    const cam = newElement("fieldset", { className: "cam" });
    const legend = newElement("legend", { innerText: name });
    const truckIcon = newElement("i", { className: "fas fa-truck" });
    const humanIcon = newElement("i", { className: "fas fa-street-view" })
    const inputIcon = newElement("i", { className: "fa-solid fa-traffic-light" });

    [legend, truckIcon, humanIcon, inputIcon].forEach(el => cam.appendChild(el));
    camerasDivEl.appendChild(cam);
}

function newElement(tagName, options = {}) {
    const el = document.createElement(tagName);
    for (const prop of Object.keys(options)) {
        el[prop] = options[prop];
    }

    return el
}

formatStatus = (status) => status ? " ready" : ""

function updateCountdown(timeLeft = 0) {
    const hours = formatNumber(Math.floor(timeLeft / 3600));
    const minutes = formatNumber(Math.floor(timeLeft / 60 - hours * 60));
    const seconds = formatNumber(timeLeft % 60);

    countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
}

formatNumber = (num) => num < 10 ? "0" + num: num
const zone = window.location.pathname.slice("/zones/".length);

const red   = "rgb(178, 49, 49)";
const green = "rgb(19, 154, 19)";

const truckIconClassName = "fas fa-truck";
const humanIconClassName = "fas fa-street-view";
const inputIconClassName = "fa-solid fa-traffic-light";

const linkClassName      = "fa-solid fa-link"
const linkSlashClassName = "fa-solid fa-link-slash"

const heartPulseClassName = "fa-solid fa-heart"
const heartCrackClassName = "fa-solid fa-heart-crack"

let timeLeft    = 0;
let cameraIDs   = [];

window.onload = startPolling;

async function get(url = "", format = "json") {
    const resp = await fetch(`${zone}/${url}`);
    if (!resp.ok) {
        throw Error(`code ${resp.status}: ${await resp.text()}`)
    }

    return format.toLowerCase() === "text" ? resp.text() : resp.json();
}

async function startPolling() {
    await render();
    const countdownEl = document.getElementById("countdown");
    const statusBtnEl = document.getElementById("status-button");
    const camerasDivEl = document.getElementById("cams");
    const webpointEl = document.getElementById("webpoint");
    const heartbeatEl = document.getElementById("heartbeat");
    statusBtnEl.addEventListener("click", handleStatusButton);

    setTimeout(async function again() {
        await updateWebpoint(webpointEl);
        await updateHeartbeat(heartbeatEl);
        await updateCameras(camerasDivEl);

        const prevTimeLeft = timeLeft;
        timeLeft = await get("countdown");
        if (timeLeft !== prevTimeLeft) {
            updateCountdown(countdownEl, timeLeft);
            updateStatusButton(statusBtnEl, timeLeft);
        }

        setTimeout(again, 1000);
    });
}

async function render() {
    const zoneEl = newElement("p", {
        id: "zone",
        innerText: await get("zone-name")
    });
    const countdownEl = updateCountdown(newElement("p", { id: "countdown" }));
    const statusBtnEl = updateStatusButton(newElement("button", {
        id: "status-button",
        innerText: "Взвесить"
    }));
    const camerasDivEl = newElement("div", { id: "cams" });

    const statusbar = newElement("div", { id: "statusbar" });
    const webpoint = newElement("i", {
        className: linkSlashClassName,
        id: "webpoint"
    });
    const heartbeat = newElement("i", {
        className: heartCrackClassName,
        id: "heartbeat"
    });
    [webpoint, heartbeat].forEach(el => statusbar.appendChild(el));

    [zoneEl, countdownEl, statusBtnEl, camerasDivEl, statusbar].forEach(el => document.body.appendChild(el));
    document.getElementById("spinner").style.display = "none";
}

function newElement(tagName, options = {}) {
    const el = document.createElement(tagName);
    for (const prop of Object.keys(options)) {
        el[prop] = options[prop];
    }

    return el
}

async function updateWebpoint(el) {
    console.log("updating webpoint...");
    const isOk = await get("webpoint");
    if (isOk) {
        el.className = linkClassName;
        el.style.color = green;
        return
    }

    el.className = linkSlashClassName;
    el.style.color = red;
}

async function updateHeartbeat(el) {
    console.log("updating heartbeat...");
    const isOk = await get("heartbeat");
    if (isOk) {
        el.className = heartPulseClassName;
        el.style.animation = "heartbeat 1s infinite";
        el.style.color = green;
        return
    }

    el.className = heartCrackClassName;
    el.style.animation = "none";
    el.style.color = red;
}

function updateCountdown(countdownEl, timeLeft = 0) {
    console.log("updating countdown...");
    const hours = formatNumber(Math.floor(timeLeft / 3600));
    const minutes = formatNumber(Math.floor(timeLeft / 60 - hours * 60));
    const seconds = formatNumber(timeLeft % 60);

    countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
    countdownEl.style.color = timeLeft ? red : green;
    return countdownEl;
}

formatNumber = (num) => num < 10 ? "0" + num: num;

function updateStatusButton(statusBtnEl, timeLeft = 0) {
    if (timeLeft) {
        statusBtnEl.disabled = true;
        statusBtnEl.style.backgroundColor = red;
        statusBtnEl.style.borderColor     = "rgb(94, 14, 14)";
        return
    }

    statusBtnEl.disabled = false;
    statusBtnEl.style.backgroundColor = green;
    statusBtnEl.style.borderColor     = "rgb(10, 78, 10)";
    return statusBtnEl;
}

async function handleStatusButton() {
    try {
        await get("button-press", "text");
        console.log("button was pressed successfully");
    } catch (err) {
        console.error(err);
    }
}

function createCamera(camerasDivEl, cameraID, states) {
    console.log(`creating camera ${cameraID}...`);
    const {name, car, human, inputs} = states;
    const camera = newElement("fieldset", { className: "cam", id: cameraID });
    const legend = newElement("legend", { innerText: name });

    const truckIcon = newElement("i", { className: truckIconClassName });
    const humanIcon = newElement("i", { className: humanIconClassName });
    const inputIcons = Object.values(inputs).map(inp => newElement("i", {
        className: inputIconClassName,
        id: inp.id
    }));

    setStatus(truckIcon, car);
    setStatus(humanIcon, human);
    inputIcons.forEach(icon => setStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state));

    [legend, truckIcon, humanIcon, ...inputIcons].forEach(el => camera.appendChild(el));
    camerasDivEl.appendChild(camera);
    return camera
}

async function updateCameras(camerasDivEl) {
    const prevCameraIDs = cameraIDs;
    cameraIDs = await get("cameras-id");

    const toRemove = prevCameraIDs.filter(prevID => !cameraIDs.find(currID => prevID === currID));
    toRemove.forEach(id => {
        const el = document.getElementById(id);
        el?.parentNode.removeChild(el);
    });

    for (const id of cameraIDs) {
        const camera = document.getElementById(id);
        const states = await get(`cameras/${id}`)
        camera ? updateCamera(camera, states) : createCamera(camerasDivEl, id, states);
    }
}

function updateCamera(camera, states) {
    console.log("updating camera states...");
    const {car, human, inputs} = states;
    for (const icon of camera.getElementsByTagName("i")) {
        if (icon.className.includes(truckIconClassName)) {
            setStatus(icon, car);
        } else if (icon.className.includes(humanIconClassName)) {
            setStatus(icon, human);
        } else if (icon.className.includes(inputIconClassName)) {
            setStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state);
        }
    }
}

function setStatus(icon, status) {
    const isReady = icon.classList.contains("ready");
    if (isReady && !status || !isReady && status) {
        icon.classList.toggle("ready");
    }
}

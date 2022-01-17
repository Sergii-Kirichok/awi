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

async function get(url = "", format = "json") {
    const resp = await fetch(`${zone}/${url}`);
    if (!resp.ok) {
        throw Error(`code ${resp.status}: ${await resp.text()}`)
    }

    return format.toLowerCase() === "text" ? resp.text() : resp.json();
}

function newElement(tagName, options = {}) {
    const el = document.createElement(tagName);
    for (const prop of Object.keys(options)) {
        el[prop] = options[prop];
    }

    return el
}

class App {
    constructor() {
        this.isHealthy = true;
        this.timeLeft  = 0;
        this.cameraIDs = [];
    }

    async render() {
        document.body.innerHTML = `
        <p id="zone">${await get("zone-name")}</p>
        <p id="countdown"></p>
        <button id="status-button">Взвесить</button>
        <div id="cameras"></div>
        <div id="statusbar">
            <i class="${linkSlashClassName}"  id="webpoint"></i>
            <i class="${heartCrackClassName}" id="heartbeat"></i>
        </div>`;

        this.countdownEl = document.getElementById("countdown");
        this.statusBtnEl = document.getElementById("status-button");
        this.camerasDivEl = document.getElementById("cameras");
        this.webpointEl = document.getElementById("webpoint");
        this.heartbeatEl = document.getElementById("heartbeat");

        this.updateCountdown();
        this.updateStatusButton();
        this.statusBtnEl.addEventListener("click", this.handleStatusButton);
    }

    spin() {
        document.body.innerHTML = `
        <div id="alert-container">
            <i class="fas fa-spinner fa-pulse" id="spinner"></i>
            <fieldset>
                <p class="alert">Нет соединения с веб-сервером</p>
            </fieldset>
        </div>`;
    }

    async update() {
        const prevTimeLeft = this.timeLeft;

        try {
            await this.updateWebpoint();
            await this.updateHeartbeat();
            await this.updateCameras();

            this.timeLeft = await get("countdown");
            if (!this.isHealthy) await this.render();
            this.isHealthy = true;
        } catch (err) {
            if (this.isHealthy) this.spin();
            this.isHealthy = false
            return
        }

        if (this.timeLeft !== prevTimeLeft) {
            this.updateCountdown(this.timeLeft);
            this.updateStatusButton(this.timeLeft);
        }
    }

    async updateWebpoint() {
        console.log("updating webpoint...");
        const isOk = await get("webpoint");
        if (isOk) {
            this.webpointEl.className   = linkClassName;
            this.webpointEl.style.color = green;
            return
        }

        this.webpointEl.className   = linkSlashClassName;
        this.webpointEl.style.color = red;
    }

    async updateHeartbeat() {
        console.log("updating heartbeat...");
        const isOk = await get("heartbeat");
        if (isOk) {
            this.heartbeatEl.className       = heartPulseClassName;
            this.heartbeatEl.style.animation = "heartbeat 1s infinite";
            this.heartbeatEl.style.color     = green;
            return
        }

        this.heartbeatEl.className       = heartCrackClassName;
        this.heartbeatEl.style.animation = "none";
        this.heartbeatEl.style.color     = red;
    }

    updateCountdown(timeLeft = 0) {
        console.log("updating countdown...");
        const hours   = this.formatNumber(Math.floor(timeLeft / 3600));
        const minutes = this.formatNumber(Math.floor(timeLeft / 60 - hours * 60));
        const seconds = this.formatNumber(timeLeft % 60);

        this.countdownEl.innerText   = `${hours}:${minutes}:${seconds}`;
        this.countdownEl.style.color = timeLeft ? red : green;
    }

    formatNumber = (num) => num < 10 ? "0" + num: num;

    updateStatusButton(timeLeft = 0) {
        if (timeLeft) {
            this.statusBtnEl.disabled              = true;
            this.statusBtnEl.style.backgroundColor = red;
            this.statusBtnEl.style.borderColor     = "rgb(94, 14, 14)";
            return
        }

        this.statusBtnEl.disabled              = false;
        this.statusBtnEl.style.backgroundColor = green;
        this.statusBtnEl.style.borderColor     = "rgb(10, 78, 10)";
    }

    async handleStatusButton() {
        try {
            await get("button-press", "text");
            console.log("button was pressed successfully");
        } catch (err) {
            console.error(err);
        }
    }

    createCamera(cameraID, states) {
        console.log(`creating camera ${cameraID}...`);
        const {name, car, human, inputs} = states;
        const camera = newElement("fieldset", { className: "camera", id: cameraID });
        const legend = newElement("legend", { innerText: name });

        const truckIcon = newElement("i", { className: truckIconClassName });
        const humanIcon = newElement("i", { className: humanIconClassName });

        let inputNum = 0;
        const inputIcons = Object.entries(inputs).map(([name, inp]) => {
            const inputEl = newElement("i", {
                className: inputIconClassName + " tooltip",
                id: inp.id
            })
            inputEl.appendChild(newElement("span", {
                className: "tooltip-text",
                // innerText: `Input: ${name}`,
                innerText: `Input: ${++inputNum}`,
            }))

            return inputEl
        });

        this.setStatus(truckIcon, car);
        this.setStatus(humanIcon, human);
        inputIcons.forEach(icon => this.setStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state));

        [legend, truckIcon, humanIcon, ...inputIcons].forEach(el => camera.appendChild(el));
        this.camerasDivEl.appendChild(camera);
        return camera
    }

    async updateCameras() {
        const prevCameraIDs = this.cameraIDs;
        this.cameraIDs = await get("cameras-id");

        const toRemove = prevCameraIDs.filter(prevID => !this.cameraIDs.find(currID => prevID === currID));
        toRemove.forEach(id => {
            const el = document.getElementById(id);
            el?.parentNode.removeChild(el);
        });

        for (const id of this.cameraIDs) {
            const camera = document.getElementById(id);
            const states = await get(`cameras/${id}`)
            camera ? this.updateCamera(camera, states) : this.createCamera(id, states);
        }
    }

    updateCamera(camera, states) {
        console.log(`updating camera ${camera.id} states...`);
        const {car, human, inputs} = states;
        for (const icon of camera.getElementsByTagName("i")) {
            if (icon.className.includes(truckIconClassName)) {
                this.setStatus(icon, car);
            } else if (icon.className.includes(humanIconClassName)) {
                this.setStatus(icon, human);
            } else if (icon.className.includes(inputIconClassName)) {
                this.setStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state);
            }
        }
    }

    setStatus(icon, status) {
        const isReady = icon.classList.contains("ready");
        if (isReady && !status || !isReady && status) {
            icon.classList.toggle("ready");
        }
    }
}

window.onload = async () => {
    const app = new App();
    await app.render();
    setTimeout(function cycle() {
        app.update();
        setTimeout(cycle, 1000);
    });
}

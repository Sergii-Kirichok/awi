const zone = window.location.pathname.slice("/zones/".length);

const truckIconClassName = "icon truck";
const humanIconClassName = "icon human";
const inputIconClassName = "icon input";

const linkSlashClassName = "icon link-slash"
const heartCrackClassName = "icon heart-crack";

const circleCheckClassName = "icon circle-check"
const circleXmarkClassName = "icon circle-xmark";

async function get(url = "", format = "json") {
    const resp = await fetch(`${zone}/${url}`);
    if (!resp.ok) {
        throw new FetchError(resp.status, await resp.text());
    }

    return format.toLowerCase() === "text" ? resp.text() : resp.json();
}

class FetchError extends Error {
    constructor(status, text) {
        super(`code ${status}: ${text}`);
        this.name = this.constructor.name;
    }
}

function newElement(tagName, options = {}) {
    const el = document.createElement(tagName);
    for (const prop of Object.keys(options)) {
        el[prop] = options[prop];
    }

    return el;
}

class App {
    constructor() {
        this.isHealthy = false;
        this.timeLeft  = 0;
        this.cameras = {};
    }

    render(name = "") {
        document.body.innerHTML = `
        <p id="zone">${name}</p>
        <p id="countdown">00:00:00</p>
        <button id="status-button" disabled>Зважити</button>
        <div id="status-bar">
            <div id="cameras"></div>
            <span id="status"></span>
        </div>`;

        this.countdownEl = document.getElementById("countdown");
        this.statusBtnEl = document.getElementById("status-button");
        this.camerasDivEl = document.getElementById("cameras");
        this.statusEl = document.getElementById("status");

        this.statusBtnEl.addEventListener("click", this.handleStatusButton.bind(this));
    }

    error(message, spinner) {
        document.body.innerHTML = `
        <div id="alert-container">
            <svg style="display: ${spinner ? "block" : "none"}" aria-hidden="true" focusable="false" class="icon spinner" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
                <path fill="rgb(35, 93, 164)" d="M96 256c0-26.5-21.5-48-48-48S0 229.5 0 256s21.5 48 48 48S96 282.5 96 256zM108.9 60.89c-26.5 0-48.01 21.49-48.01 47.99S82.39 156.9 108.9 156.9s47.99-21.51 47.99-48.01S135.4 60.89 108.9 60.89zM108.9 355.1c-26.5 0-48.01 21.51-48.01 48.01S82.39 451.1 108.9 451.1s47.99-21.49 47.99-47.99S135.4 355.1 108.9 355.1zM256 416c-26.5 0-48 21.5-48 48S229.5 512 256 512s48-21.5 48-48S282.5 416 256 416zM464 208C437.5 208 416 229.5 416 256s21.5 48 48 48S512 282.5 512 256S490.5 208 464 208zM403.1 355.1c-26.5 0-47.99 21.51-47.99 48.01S376.6 451.1 403.1 451.1s48.01-21.49 48.01-47.99S429.6 355.1 403.1 355.1zM256 0C229.5 0 208 21.5 208 48S229.5 96 256 96s48-21.5 48-48S282.5 0 256 0z"></path>
            </svg>
            <fieldset>
                <p class="alert">${message}</p>
            </fieldset>
        </div>`;
    }

    async update() {
        const prevTimeLeft = this.timeLeft;

        try {
            const { name, heartbeat, webpoint, timeLeft, cameras } = await get("data");
            if (!this.isHealthy) this.render(name);
            this.isHealthy = true;

            this.updateWebpoint(webpoint, heartbeat);
            this.updateCameras(cameras);

            this.timeLeft = timeLeft;
        } catch (err) {
            this.isHealthy = false
            console.error(`updating error: ${err.message}`);
            if (err instanceof FetchError) {
                this.error(err.message, false)
                return
            }

            this.error("Нет соединения с веб-сервером", true)
        }

        if (this.timeLeft !== prevTimeLeft) {
            this.updateCountdown(this.timeLeft);
            this.updateStatusButton(this.timeLeft);
        }
    }

    updateWebpoint(webpoint, heartbeat) {
        // console.log("updating webpoint...");
        if (!webpoint) {
            this.statusEl.className = linkSlashClassName;
            this.statusEl.style.display = "block";
            return
        }

        // console.log("updating heartbeat...");
        if (!heartbeat) {
            this.statusEl.className = heartCrackClassName;
            this.statusEl.style.display = "block";
            return
        }

        this.statusEl.className = "";
        this.statusEl.style.display = "none";
    }

    updateCountdown(timeLeft = 0) {
        // console.log("updating countdown...");
        const hours   = this.formatNumber(Math.floor(timeLeft / 3600));
        const minutes = this.formatNumber(Math.floor(timeLeft / 60 - hours * 60));
        const seconds = this.formatNumber(timeLeft % 60);

        this.countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
        this.setStatus(this.countdownEl, !timeLeft);
    }

    formatNumber = (num) => num < 10 ? "0" + num: num;

    updateStatusButton(timeLeft = 0) {
        this.setStatus(this.statusBtnEl, !timeLeft);
        this.statusBtnEl.disabled = !!timeLeft;
    }

    async handleStatusButton() {
        try {
            await get("button-press", "text");
            console.log("button was pressed successfully");
            this.statusEl.className = circleCheckClassName;
            this.statusEl.style.display = "block";
            setTimeout(() => this.statusEl.style.display = "none", 500);
        } catch (err) {
            console.error(err);
            this.statusEl.className = circleXmarkClassName;
            this.statusEl.style.display = "block";
        }
    }

    createCamera(cameraID, states) {
        console.log(`creating camera ${cameraID}...`);
        const {name, car, human, inputs} = states;
        const camera = newElement("fieldset", { className: "camera", id: cameraID });
        const legend = newElement("legend", { innerText: name });

        const truckIcon = newElement("span", { className: truckIconClassName });
        const humanIcon = newElement("span", { className: humanIconClassName });

        let inputNum = 0;
        const inputIcons = Object.entries(inputs).map(([name, inp]) => {
            const inputEl = newElement("span", {
                className: inputIconClassName + " tooltip",
                id: inp.id
            })
            inputEl.appendChild(newElement("span", {
                className: "tooltip-text",
                // innerText: `Input: ${name}`,
                innerText: `Input: ${++inputNum}`,
            }))

            return inputEl;
        });

        this.setStatus(truckIcon, car);
        this.setStatus(humanIcon, human);
        inputIcons.forEach(icon => this.setStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state));

        [legend, truckIcon, humanIcon, ...inputIcons].forEach(el => camera.appendChild(el));
        this.camerasDivEl.appendChild(camera);
        return camera;
    }

    updateCameras(cameras) {
        const prevCameras = this.cameras;
        this.cameras = cameras;

        const toRemove = Object.keys(prevCameras).filter(prevID => !Object.keys(this.cameras).find(currID => prevID === currID));
        toRemove.forEach(id => {
            const el = document.getElementById(id);
            el?.parentNode.removeChild(el);
        });

        Object.entries(cameras).forEach(([id, states]) => {
            const camera = document.getElementById(id);
            camera ? this.updateCamera(camera, states) : this.createCamera(id, states);
        })
    }

    updateCamera(camera, states) {
        // console.log(`updating camera ${camera.id} states...`);
        const {car, human, inputs} = states;
        for (const icon of camera.getElementsByTagName("span")) { // todo: pay attention
            if (icon.className.includes(truckIconClassName)) {
                this.setStatus(icon, car);
            } else if (icon.className.includes(humanIconClassName)) {
                this.setStatus(icon, human);
            } else if (icon.className.includes(inputIconClassName)) {
                this.setStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state);
            }
        }
    }

    setStatus(element, status) {
        const isReady = element.classList.contains("ready");
        if (isReady && !status || !isReady && status) {
            element.classList.toggle("ready");
        }
    }
}

window.onload = async () => {
    const app = new App();
    setTimeout(function cycle() {
        app.update();
        setTimeout(cycle, 1000);
    });
};

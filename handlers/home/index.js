const green       = "rgb(19, 154, 19)";
const greenBorder = "rgb(10, 78, 10)";
const red         = "rgb(178, 49, 49)";
const rebBorder   = "rgb(94, 14, 14)";

const h = 0;
const m = 0;
const s = 5;

let time;

const countdownEl = document.getElementById("countdown");
const statusBtnEl = document.getElementById("status-button");

statusBtnEl.addEventListener("click", () => { if (!time) recover() })

function recover() {
    statusBtnEl.disabled = true;
    statusBtnEl.style.backgroundColor = red;
    statusBtnEl.style.borderColor = rebBorder;
    time = (h * 60 + m) * 60 + s;
    countdown();
}

window.onload = () => recover();

function countdown() {
    setTimeout(function again() {
        updateCountdown()
        if (!time) {
            statusBtnEl.style.backgroundColor = green;
            statusBtnEl.style.borderColor = greenBorder;
            statusBtnEl.disabled = false;
            return
        }

        setTimeout(again, 1000)
    }, 1000)
}

function updateCountdown() {
    time--;

    const hours = decrease(Math.floor(time / 3600));
    const minutes = decrease(Math.floor(time / 60 - hours * 60));
    const seconds = decrease(time % 60);

    countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
}

decrease = (num) => num < 10 ? "0" + num: num
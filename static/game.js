var ip = window.location.href.split("://")[1].slice(0, -1);
console.log(ip);
const socket = io("http://" + ip);

var choice;
const p1Span = document.getElementById('score-p1');
const p2Span = document.getElementById('score-p2');
const p1ChoiceImg = document.getElementById('p1-choice');
const p2ChoiceImg = document.getElementById('p2-choice');
const paperChoice = document.getElementById('paper');
const rockChoice = document.getElementById('rock');
const scissorsChoice = document.getElementById('scissors');
const spanMsg = document.getElementById('message');
const choiceList = document.getElementById('choice-list');
var p1Choice;
var p2Choice;

socket.on("connect", () => {
    console.log("socket connected");
})

socket.on("gameReset", function (response) {
    console.log("game reset");
    location.reload();
})

socket.on("gameStart", function (response) {
    console.log("game start");
    p1Span.innerHTML = "0";
    p2Span.innerHTML = "0";
    choiceList.style.visibility = "unset";
    spanMsg.innerHTML = "Make your choice";
})

socket.on("makeMove", function (response) {
    spanMsg.innerHTML = "Making move...";
    
    response = response.split("#");
    p1Choice = response[0];
    p2Choice = response[1];
    p1Score = response[2];
    p2Score = response[3];

    update_ui(response[4]);
})

socket.on('ready', function (response) {
    setTimeout(function () {
        console.log(response);

        document.getElementById("intro-menu").style.display = "none";
        if (response == "") {
            document.getElementById("game-full").style.display = "block";
        }
        else {
            document.getElementById("game").style.display = "block";
            document.getElementById("current-player").innerHTML = response;
            document.getElementById(response.toLowerCase() + "-choice").classList.add("highlight")
        }

    }, 500);
});

rockChoice.addEventListener('click', function () {
    choice = 'rock';
    play();
});

paperChoice.addEventListener('click', function () {
    choice = 'paper';
    play();
});

scissorsChoice.addEventListener('click', function () {
    choice = 'scissors';
    play();
});

function play() {
    choiceList.style.visibility = "hidden";
    spanMsg.innerHTML = "Waiting for the other player";
    socket.emit("play", choice);
}

function update_ui(msg) {
    p1ChoiceImg.src = 'images/rock.png';
    p2ChoiceImg.src = 'images/rock.png';
    choiceList.style.visibility = 'hidden';
    p2ChoiceImg.style.animation = 'p1-shake 2s ease';
    p1ChoiceImg.style.animation = 'p2-shake 2s ease';

    setTimeout(function () {
        p1Span.innerText = p1Score;
        p2Span.innerText = p2Score;
        p2ChoiceImg.src = "images/" + p2Choice + ".png";
        p1ChoiceImg.src = "images/" + p1Choice + ".png";
        spanMsg.innerHTML = msg;
        choiceList.style.visibility = 'unset';
        p2ChoiceImg.style.animation = '';
        p1ChoiceImg.style.animation = '';
    }, 2000)
}
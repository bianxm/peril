# Peril: Crowdsourced Jeopardy (Back end)

For when you want to play Jeopardy with your friends, with each friend having their own category that they write the questions for. When their category is selected, that friend is the 'quizmaster' for that question, and will read out their question and set when buzzers are ready to be pressed, then mark the answer right or wrong. 

Inspired by Jackbox games, the aim is to have a main screen that will be streamed to all participants. The game will be started from the screen, then participants can join on their phones using the game code, then input their screen names and category + questions. The main screen will trigger the game starting. During the game, the main screen will display the questions and the overall scores, while each player's phone will function as their Jeopardy buzzer.

## To do
- [x] Set up Websocket communications boilerplate (greatly indebted to [dhij/go-next-ts_chat](https://github.com/dhij/go-next-ts_chat))
- [ ] Game logic
    - [x] Switching between game stages
    - [ ] Keeping track of scores
    - [ ] Input validation
- [ ] Game communication
    - [ ] Send only relevant game state info to player clients
    - [ ] Send only relevant game state info to 'screen'
- [ ] Pregame set-up
    - [ ] Username input
    - [ ] Category and Question input
- [ ] Front end
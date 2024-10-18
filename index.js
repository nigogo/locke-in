/**
 * @param text {string}
 * @param date {Date}
 * @returns {{type: string, date: Date, text: string}}
 */
function setGoal(text, date) {
  return {
    type: 'SET_GOAL',
    text,
    date
  }
}

const result = setGoal('Finish LockeIn #nobuild setup', new Date());
console.log(result);
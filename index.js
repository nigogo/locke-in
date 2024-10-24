import Koa from "koa";
import Router from "@koa/router";
import Pug from "koa-pug";
import serve from "koa-static";
import * as path from "node:path";

const app = new Koa();
const router = new Router();
const pug = new Pug({
    viewPath: './src/views',
    basedir: './src/views',
    app: app
});

app.use(serve( './static'));

router.get("/", async (ctx) => {
    ctx.body = await pug.render('index');
});

router.post("/goal", async (ctx) => {
    ctx.body = await pug.render('goal', {
        goal: "Finish LockeIn #nobuild setup",
        datetime: new Date().setSeconds(new Date().getSeconds() + 5)
    });
});

app
    .use(router.routes())
    .use(router.allowedMethods())
    .listen(3000);

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
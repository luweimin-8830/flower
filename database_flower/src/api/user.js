//模版页面,新接口页面直接复制此页
import { Router } from "express";
import { ok, fail } from "../response.js";
import { User } from "../../db.js";

const router = Router();

router.get("/add", async (req, res) => {
    try {
        const OPENID = req.headers['x-wx-openid']
        const addUser = await User.create({
            openId:OPENID,
            lastLogin: new Date()
        })
        res.json(ok(addUser))
    } catch (e) { console.log(e) }
})

export default router;
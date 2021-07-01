import axios from "axios";
// const baseURL = "http://127.0.0.1:1232"
const baseURL = ""

export default axios.create({
  baseURL,
  headers: {
    "Content-Type": "application/json",
    "charset": "utf-8"
  }
})
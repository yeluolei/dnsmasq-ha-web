import request from "./request";

export default {
  getHosts() {
    return request.get("/hosts");
  },

  createHost(host) {
    return request.post("/hosts", host);
  },

  updateHost(host) {
    return request.put(`/hosts/${host.id}`, host);
  },

  deleteHost(id) {
    return request.delete(`/hosts/${id}`);
  },

  applyChange()  {
    return request.post("/generate")
  }
}
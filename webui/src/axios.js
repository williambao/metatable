import axios from 'axios'
import router from './router/index'
import store from './store/index'
import * as types from './store/types'

// console.log(process.env.BASE_URL)
axios.defaults.baseURL = process.env.BASE_URL;
// axios.defaults.baseURL = 'http://192.168.199.130:3000';
// axios.defaults.timeout = 5000;
axios.defaults.headers.post['Content-Type'] = 'application/json';

const instance = axios.create();
instance.defaults.headers.post['Content-Type'] = 'application/json';

axios.interceptors.request.use = instance.interceptors.request.use;
instance.interceptors.request.use(config => {
    if (localStorage.getItem('token')) {
        config.headers.Authorization = `Bearer ${localStorage.getItem('token')}`
            .replace(/(^\")|(\"$)/g, '')
    }
    return config
}, err => {
    return Promise.reject(err)
})

// axios拦截响应
instance.interceptors.response.use(response => {
    return response
}, err => {
    return Promise.reject(err)
})

export default {
    register(data) {
        return instance.post('/account/register', data);
    },
    login(data) {
        return instance.post('/account/login', data);
    },
    mySelf() {
        return instance.get('/user');
    },
    getUserTables(userId) {
        return instance.get('/user/' + userId + '/tables');
    },
    updateUserInfo(id, data) {
        return instance.put('/user/' + id, data);
    },
    updatePassword(id, data) {
        return instance.post('/user/' + id + '/password', data);
    },
    getFileUploadToken() {
        return instance.get('/files/token');
    },
    getTable(id) {
        return instance.get('/tables/' + id);
    },

    createTable(table) {
        return instance.post("/tables", table);
    },
    updateTable(tableId, data) {
        return instance.put("/tables/" + tableId, data);
    },
    getTableUsers(id) {
        return instance.get("/tables/" + id + "/users");
    },
    deleteTable(id) {
        return instance.delete("/tables/" + id);
    },
    addTableUser(tableId, username) {
        return instance.post("/tables/" + tableId + "/users", { "username": username });
    },
    updateTableUser(tableId, id, data) {
        return instance.put("/tables/" + tableId + "/users/" + id, data);
    },
    deleteTableUser(tableId, id) {
        return instance.delete("/tables/" + tableId + "/users/" + id);
    },
    // 
    createColumn(data) {
        return instance.post('/tables/' + data.tableId + '/columns', data)
    },
    updateColumn(data) {
        return instance.put('/tables/' + data.tableId + '/columns/' + data.id, data)
    },
    deleteColumn(tableId, columnId) {
        return instance.delete('/tables/' + tableId + '/columns/' + columnId);
    },
    updateColumnOrders(tableId, data) {
        return instance.put('/tables/' + tableId + '/column-orders', data);
    },
    addRecord(data) {
        return instance.post('/tables/' + data.tableId + '/data', data);
    },
    deleteRecord(tableId, recordId) {
        return instance.delete('/tables/' + tableId + '/data/' + recordId);
    },
    deleteRecords(tableId, recordIds) {
        return instance.delete('/tables/' + tableId + '/data?ids=' + recordIds);
    },
    getRecords(tableId, query) {
        var url = '/tables/' + tableId + '/data?_=_';
        if (query.searchText) {
            url += ("&query=" + query.searchText);
        }
        // debugger;
        if (query.sort) {
            url += ("&sort=" + query.sort);
        }
        if (query.hide_columns) {
            url += ("&hide_columns=" + query.hide_columns);
        }
        return instance.get(url);
    },
    getViewRecords(tableId, viewId, query) {
        var url = '/tables/' + tableId + '/views/' + viewId + '/data?_=_';
        if (query.searchText) {
            url += ("&query=" + query.searchText);
        }
        return instance.get(url);
    },
    editRecordCell(tableId, recordId, cells) {
        return instance.patch('/tables/' + tableId + '/data/' + recordId, cells);
    },

    updateTableView(tableId, viewId, view) {
        return instance.put('/tables/' + tableId + '/views/' + viewId, view);
    },


    getConditions() {
        return instance.get('/dict/conditions');
    },

    getTemplates(page, pagesize) {
        return instance.get("/templates?page=" + page + "&pagesize=" + pagesize);
    },
    createTemplate(table) {
        return instance.post("/templates", table);
    },
    deleteTemplate(id) {
        return instance.delete("/templates/" + id);
    },
    getTemplate(id) {
        return instance.get("/templates/" + id);
    },
    updateTemplate(id, data) {
        return instance.put("/templates/" + id, data);
    },
    updateTemplateColumn(id, columnId, data) {
        return instance.put("/templates/" + id + '/columns/' + columnId, data);
    },
    addTemplateColumn(id, data) {
        return instance.post("/templates/" + id + '/columns', data);
    },
    getTemplateColumns(id) {
        return instance.get("/templates/" + id + "/columns");
    },
    deleteTemplateColumn(id, columnId) {
        return instance.delete("/templates/" + id + "/columns/" + columnId);
    },
    updateTemplateColumnOrders(id, data) {
        return instance.put('/templates/' + id + '/column-orders', data);
    },

    addOrganization(data) {
        return instance.post("/organizations", data);
    },
    getOrganizations(userId) {
        return instance.get("/organizations");
    },
    updateOrganization(id, data) {
        return instance.put("/organizations/" + id, data);
    },
    deleteOrganization(id) {
        return instance.delete("/organizations/" + id);
    },
    getOrganizationUsers(id) {
        return instance.get("/organizations/" + id + "/members");
    },
    addOrganizationUser(id, data) {
        return instance.post("/organizations/" + id + "/members", data);
    },
    updateOrganizationUser(id, userId, data) {
        return instance.put("/organizations/" + id + "/members/" + userId, data);
    },
    deleteOrganizationUser(id, userId) {
        return instance.delete("/organizations/" + id + "/members/" + userId);
    },
}
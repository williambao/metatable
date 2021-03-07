import * as types from './types'

const mutations = {
    [types.IS_LOADING]: (state, data) => {
        state.isLoading = data;
    },
    [types.LOGIN]: (state, data) => {
        localStorage.setItem('token', data.token);
        localStorage.setItem('user', JSON.stringify(data));
        state.user = data;
        state.token = data.token;
    },
    [types.LOGOUT]: (state) => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        state.user = {};
        state.token = null;
    },
    [types.SHOW_NAV_TABLES]: (state, data) => {
        state.showNavTables = data;
    },
    [types.USER]: (state, data) => {
        localStorage.setItem('user', JSON.stringify(data));
        state.user = data;
    },
    [types.SHOW_NAV]: (state, data) => {
        state.showNav = data;
    },
    [types.TABLES]: (state, data) => {
        state.tables = data;
    },
    [types.COLUMNS]: (state, data) => {
        state.columns = data;
    },
    [types.CONDITIONS]: (state, data) => {
        state.conditions = data;
    },
    [types.VIEWS]: (state, data) => {
        state.views = data;
    },
    [types.SELECT_TABLE]: (state, data) => {
        state.table = data;
    },
    [types.TABLE_PERMISSION]: (state, data) => {
        state.permission = data;
    },
    [types.TABLE_VIEW]: (state, data) => {
        state.view = data;
    },
    [types.RECORDS]: (state, data) => {
        state.records = data;
    },

    [types.FILE_UPLOAD_TOKEN]: (state, data) => {
        state.fileUploadToken = data;
    },


    [types.TABLE_SEARCH_TEXT]: (state, data) => {
        state.tableSearchText = data;
    },
    [types.TABLE_FILTERS]: (state, data) => {
        state.tableFilters = data;
    },
    [types.TABLE_HIDDEN_COLUMNS]: (state, data) => {
        state.tableHiddenColumns = data;
    },
    [types.TABLE_SORTS]: (state, data) => {
        state.tableFilters = data;
    },
    [types.ORGANIZATIONS]: (state, data) => {
        state.organizations = data;
    },
    [types.MEMBERS]: (state, data) => {
        state.members = data;
    },

}

export default mutations
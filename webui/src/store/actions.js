import Vue from 'vue'
import router from '../router/index'
import * as types from './types'
import api from '../axios'

export default {
    IsLoading({ commit }, data) {
        commit(types.IS_LOADING, data);
    },
    UserLogin({ commit }, data) {
        commit(types.LOGIN, data)
    },

    UserLogout({ commit }) {
        commit(types.LOGOUT)
    },

    SetUser({ commit }, data) {
        commit(types.USER, data)
    },

    ShowNav({ commit }, data) {
        commit(types.SHOW_NAV, data)
    },

    SetTables({ commit }, data) {
        commit(types.TABLES, data);
    },

    SetViews({ commit }, data) {
        commit(types.VIEWS, data);
    },

    SetTable({ commit }, data) {
        commit(types.SELECT_TABLE, data);
    },

    SetColumns({ commit }, data) {
        commit(types.COLUMNS, data);
    },

    SetRecords({ commit }, data) {
        commit(types.RECORDS, data);
    },

    SetView({ commit }, data) {
        commit(types.TABLE_VIEW, data);
    },

    SetTablePermission({ commit }, data) {
        commit(types.TABLE_PERMISSION, data);
    },

    SetFileUploadToken({ commit }, data) {
        commit(types.FILE_UPLOAD_TOKEN, data);
    },


    SetTableSearchText({ commit }, data) {
        commit(types.TABLE_SEARCH_TEXT, data);
    },
    SetTableSorts({ commit }, data) {
        commit(types.TABLE_SORTS, data);
    },
    SetTableFilters({ commit }, data) {
        commit(types.TABLE_FILTERS, data);
    },
    SetTableHiddenColumns({ commit }, data) {
        commit(types.TABLE_HIDDEN_COLUMNS, data);
    },
    GetOrganizations({ commit }, data) {
        api.getOrganizations().then(res => {
            commit(types.ORGANIZATIONS, res.data);
        });
    },

    GetOrganizationMembers({ commit }, data) {
        api.getOrganizationUsers(data).then(res => {
            commit(types.MEMBERS, res.data);
        });
    }

}
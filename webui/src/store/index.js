import Vue from 'vue'
import Vuex from 'vuex'
import mutations from './mutations'
import actions from './actions'

Vue.use(Vuex);

const state = {
    'token': null,
    'activeName': 'first',
    'showNav': true,
    'user': {},
    'tables': [],
    'columns': [],
    'views': [],
    'conditions': {},
    'showNavTables': false,
    'table': null,
    'view': { hide_columns: [], sorts: [], filters: [] },
    'records': [],
    'members': [],
    'permission': {},
    'isLoading': false,
    'fileUploadToken': { 'token': "", 'upload_url': "" },
    'tableHiddenColumns': [],
    'tableFilters': [],
    'tableSorts': [],
    'organizations': [],
    'tableSearchText': '',
}

const getters = {
    visiableColumns: state => {
        if (!state.columns) return [];
        if (!state.view || !state.view.hide_columns) return state.columns;
        return state.columns.filter(col => state.view.hide_columns.indexOf(col.id) == -1);
    },
    hide_columns: state => {
        return state.view.hide_columns || [];
    },
    sorts: state => {
        return state.view.sorts || [];
    },
    filters: state => {
        return state.view.filters || [];
    }
}

export default new Vuex.Store({
    state,
    getters,
    mutations,
    actions,
})
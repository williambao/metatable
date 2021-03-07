<template>
    <div class="toolbar float-container">
        <div class="action item" @click="refresh()">
            <Badge >
                <Icon type="md-refresh" />
            </Badge>
            <div class="title">刷新</div>
        </div>     
        <!-- <div class="action item" :disabled="!selectedRows.length">
            <Badge >
                <Icon type="ios-compose" size="16"></Icon>
            </Badge>
            <div class="title">编辑</div>
        </div>   --> 
        <div class="action item" v-if="selectedRows.length"  @click="openDeleteDialog()">
            <Badge >
                <Icon type="md-trash" />
            </Badge>
            <div class="title">删除</div>
        </div>    
        <hr />
        <div class="action " style="height:100%;">
            <Input size="small"  placeholder="搜索..." class="search-box" v-model="searchText" @input="onInput"></Input>
        </div>
        
        <div class="title-box" v-if="table" @dblclick="showEditTable">
            <span v-if="table.organization">{{table.organization.name}} / </span>
            {{table.name}}
        </div>
        
        <div class="right">
            <div class="action item" @click="toggleFilters()" v-if="false">
                <Badge :dot="filters.length>0">
                    <Icon type="md-funnel" />
                </Badge>
                <div class="title">筛选</div>
            </div>
            <div class="action item" @click="toggleSorts()">
                <Badge :dot="sorts.length>0">
                    <Icon type="md-repeat" />
                </Badge>
                <div class="title">排序</div>
            </div>
            <div class="action item" @click="showHiddenColumnDialog=!showHiddenColumnDialog">
                <Badge :dot="hide_columns.length>0">
                    <Icon type="md-eye-off" />
                </Badge>
                <div class="title">隐藏列</div>

            </div>
            <hr />
            <div class="action item" @click="showColumnSortDialog=!showColumnSortDialog">
                <Badge >
                    <Icon type="ios-options" />
                </Badge>
                <div class="title">列顺序</div>

            </div>            
            <hr />
            <!--
            <div class="action item " disabled="disabled" @click="showNotImplementWarning()">
                <Badge >
                    <Icon type="upload" size="18"></Icon>
                </Badge>
                <div class="title">导入</div>
            </div>       
            <div class="action item " disabled="disabled" @click="showNotImplementWarning()">
                <Badge >
                    <Icon type="archive" size="18"></Icon>
                </Badge>
                <div class="title">导出</div>
            </div>       
            -->
            <div class="action item " @click="openTableUsers()">
                <Tooltip placement="bottom-end">
                    <Badge >
                        <Icon type="md-lock" />
                    </Badge>
                    <div class="title">权限</div>
                    <div slot="content">
                        <p>设置表格成员</p>
                    </div>
                </Tooltip>
            </div>
        </div>
        <Modal v-model="showColumnSortDialog" title="列显示顺序" width="300" >
            <div class="dropdown-pop-wrapper">
                <draggable v-model="columns" :options="{group:'people'}" @end="updateColumnOrders">
                    <div v-for="column in columns" class="draggable-item" :key="'sort_'+column.id">
                        <Icon type="images" v-if="column.type=='image'"></Icon>
                        <Icon type="ios-grid-view-outline" v-if="column.type=='number'"></Icon>
                        <Icon type="android-create" v-if="column.type=='string'"></Icon>
                        <Icon type="navicon" v-if="column.type=='text'"></Icon>
                        <Icon type="ios-calendar-outline" v-if="column.type=='datetime'"></Icon>
                        <Icon type="ios-checkmark-outline" v-if="column.type=='boolean'"></Icon>
                        <Icon type="chevron-down" v-if="column.type=='selection'"></Icon>
                        <Icon type="ios-information-outline" v-if="column.type=='system'"></Icon>                    
                        
                        <strong v-if="column.is_primary">{{column.name}}</strong>
                        <span v-if="!column.is_primary">{{column.name}}</span>

                        <Icon type="navicon" class="right text-gray" size="18" style="margin-top:8px;"></Icon>
                    </div>
                </draggable>
            </div>
            <div slot="footer" style="text-align:left;color:#aaa;">上下拖动列可排序显示位置</div>
        </Modal>
        <Modal v-model="showHiddenColumnDialog" title="隐藏列" width="300" >
            <div class="dropdown-pop-wrapper hide-columns">
                <Checkbox-group v-model="hide_columns" >
                    <Checkbox v-for="column in columns" :label="column.id" :key="'hide_'+column.id" > {{column.name}}</Checkbox>
                </Checkbox-group>

            </div>
            <div slot="footer" style="text-align:left;color:#aaa;">请选上不想显示的列</div>
        </Modal>
        <Modal  v-model="showSortsDialog" title="查询排序" width="300" @on-ok="updateSorts" @on-cancel="resetSorts">
            <div class="sorts">
                <div>
                    <div class="row-item" v-for="(sort, index) in sorts" :key="'x_'+sort.id">
                        <Select  v-model="sort.id">
                            <Option v-for="column in columns" :value="column.id" :key="'sort2_'+column.id">{{ column.name }}</Option>
                        </Select>
                        <Checkbox v-model="sort.is_desc" style="margin-left:5px;width:50px;">倒序</Checkbox>
                        <Icon type="ios-minus-outline" size="24" style="vertical-align:middle;cursor:pointer;" @click.native="removeSort(index)"></Icon>
                    </div>
                </div>
                <a type="text" class="add-button" @click.stop="addSort">增加排序</a>
            </div>
            
        </Modal>
        <Modal v-model="showFilterDialog" title="筛选数据" width="400" @on-ok="updateFilters" @on-cancel="resetFilters" >
            <div>
                <span>筛选出符合下面</span>
                <Select  v-model="filterType"  size="small" style="width:70px;">
                    <Option value="and" >所有</Option>
                    <Option value="or" >任一</Option>
                </Select>
                <span>条件的数据</span>
            </div>
            <div style="margin-top:20px;">
                <div class="row-item" v-for="(filter, index) in filters" :key="'filter_' + filter.id">
                    <Select v-model="filter.column_id" style="width:120px;" @on-change="setFilterTypes" :label-in-value="true">
                        <Option v-for="column in columns" :value="column.id" :label="column.name" :key="'option_' + column.id">
                            <span>
                                <Icon type="images" v-if="column.type=='image'"></Icon>
                                <Icon type="ios-grid-view-outline" v-if="column.type=='number'"></Icon>
                                <Icon type="android-create" v-if="column.type=='string'"></Icon>
                                <Icon type="navicon" v-if="column.type=='text'"></Icon>
                                <Icon type="ios-calendar-outline" v-if="column.type=='datetime'"></Icon>
                                <Icon type="ios-checkmark-outline" v-if="column.type=='boolean'"></Icon>
                                <Icon type="chevron-down" v-if="column.type=='selection'"></Icon>
                                <Icon type="ios-information-outline" v-if="column.type=='system'"></Icon>
                            </span>
                            <span>{{column.name}}</span>
                        </Option>
                    </Select>
                    <Select  v-model="filter.operator" style="width:70px;">
                        <Option v-for="condition in filter.conditions" :value="condition.type" :key="'condition_'+condition.type">{{ condition.name }}</Option>
                    </Select>
                    <Input v-if="filter.conditions.filter(t => t.type==filter.operator).length && filter.conditions.filter(t => t.type==filter.operator)[0].value" v-model="filter.value" placeholder="请输入..." style="width: 120px"></Input>
                    <Icon type="ios-minus-outline" size="24" style="vertical-align:middle;cursor:pointer;" @click.native="removeFilter(index)"></Icon>
                    
                </div>
            </div>
            <a type="text" class="add-button" @click.stop="addFilter">增加条件</a>
        </Modal>
        <Modal
            v-model="showUsersDialog"
            title="数据权限设置"
            width="650"
            >
            <div style="max-height:350px">
            <div style="padding:0 0 10px 0;" >
                <span v-if="table && !table.organization_id">
                    <Input v-model="username" placeholder="用户名" style="width:150px;"></Input>
                    <Button @click.native="addTableUser">增加用户</Button>
                </span>
                <span v-else>
                    <Select v-model="username" style="width:150px;" clearable placeholder="请选择用户">
                        <Option v-for="mm in members" :value="mm.username" :key="mm.id">{{mm.nickname}}</Option>
                    </Select>
                    <Button @click.native="addTableUser">增加用户</Button>
                </span>
            </div>
                <table class="table user-table" v-if="tableUsers.length">
                    <thead>
                        <tr>
                            <th >用户昵称</th>
                            <!-- <th>用户名</th> -->
                            <!-- <th>新增权限</th> -->
                            <th>查看数据权限</th>
                            <th>编辑数据权限</th>
                            <!-- <th>删除权限</th> -->
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(u,index) in tableUsers" :key="u.id">
                            <td>{{u.nickname}}</td>
                            <!-- <td>{{u.username}}</td> -->
                            <!-- <td><Checkbox v-model="u.is_allow_insert_data" @on-change="updateTableUser(u, 'is_allow_insert_data')"></Checkbox></td> -->
                            <td>
                                <Select v-model="u.view_data" size="small" @on-change="updateTableUser(u, 'view_data')">
                                    <Option v-for="perm in permissions" :value="perm.value" :key="perm.value">{{perm.text}}</Option>
                                </Select>   
                            </td>
                            <td>
                                <Select v-model="u.edit_data" size="small"  @on-change="updateTableUser(u, 'edit_data')">
                                    <Option v-for="perm in permissions" :value="perm.value" :key="perm.value">{{perm.text}}</Option>
                                </Select>   
                            </td>
                            <!-- <td>
                                <Select v-model="u.delete_data" size="small" @on-change="updateTableUser(u, 'delete_data')">
                                    <Option v-for="perm in permissions" :value="perm.value" :key="perm.value">{{perm.text}}</Option>
                                </Select>   
                            </td> -->
                            <td>
                                <Poptip
                                    confirm
                                    title="您确认删除此用户吗？"
                                    v-if="!u.is_table_admin"
                                    @on-ok="deleteTableUser(u)">
                                    <Button size="small">删除</Button>
                                </Poptip>
                                <span v-if="u.is_table_admin" class="text-success">管理员</span>

                            </td>
                        </tr>
                    </tbody>
                </table>
                <div v-if="!tableUsers.length">无用户</div>
            </div>
            <div slot="footer"></div>
        </Modal>
        <Modal v-model="showEditTableDialog" title="修改表格名称" @on-ok="updateTableName">
            <Input v-model="tableName" placeholder="表格名称"/>
        </Modal>
    </div>
    
</template>
<script>
import debounce from '../../util/debounce'
import * as types from '../../store/types'
import api from '../../axios'
import draggable from 'vuedraggable'
export default {
    props: ['table'],
    components: {
        draggable,
    },
    data() {
        return {
            searchText: '',
            tableName: '',
            showHiddenColumnDialog: false,
            showColumnSortDialog: false,
            showEditTableDialog: false,
            showSortsDialog: false,
            showFilterDialog: false,
            showUsersDialog: false,
            tableUserLoading: false,
            editingSort: {},
            sorts: [],
            filterType: 'and',
            filters: [],
            columnConditions: [],
            username: "",
            tableUsers: [],
            permissions: [
                {text: "全部数据",value:"all"},
                // {text: "本组人创建数据",value:"team"},
                {text: "本人创建数据",value:"user"},
                {text: "无权限",value:"none"},
            ]
        }
    },
    computed: {
        members() {
            return this.$store.state.members;
        },
        view() {
            return this.$store.state.view;
        },
        conditions() {
            return this.$store.state.conditions;
        },
        columns: {
            get() {
                return this.$store.state.columns
            },
            set(value) {
                this.$store.commit('columns', value)
            }
        },
        records() {
            return this.$store.state.records;
        },
        
        selectedRows() {
            return this.$store.state.records.filter(row => row.is_selected);
        },
        hide_columns: {
            get() {
                return this.$store.getters.hide_columns;
            },
            set(value) {
                let obj = this.view;
                obj.hide_columns = value;
                this.$store.dispatch('SetView', obj);
                this.updateHiddenColumns();
            }
        },
        view_sorts() {
            return this.$store.getters.sorts; 
        },
        view_filters() {
            return this.$store.getters.filters; 
        },
    },
    methods: {
        showEditTable() {
            this.tableName = this.table.name;
            this.showEditTableDialog = true;
        },
        openTableUsers() {
            if(!this.showUsersDialog) {
                this.getTableUsers();
            }
            this.showUsersDialog = !this.showUsersDialog;
        },
        updateTableName() {
            if(!this.tableName) {
                this.$Message.warning("请输入表格名称");
                return;
            }
            api.updateTable(this.table.id, {name: this.tableName}).then(res => {
                this.$Message.success("保存成功");
                this.$emit("refresh");
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        addTableUser() {
            if(!this.username.length) {
                this.$Message.warning("请输入要增加的用户名");
                return;
            }
            let self = this;
            api.addTableUser(this.table.id, this.username).then(res => {
                self.getTableUsers();
                this.$Message.success("新增成功");
                self.username = '';
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        deleteTableUser(u) {
            let self = this;
            api.deleteTableUser(this.table.id, u.id).then(res => {
                self.getTableUsers();
                this.$Message.success("删除成功");
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        updateTableUser(u, key) {
            let self = this;
            var data = {};
            data[key] = u[key];
            api.updateTableUser(this.table.id, u.id, data).then(res => {
                this.$Message.success("更新成功");
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        getTableUsers() {
            if(!this.table) return;
            let self = this;

            this.tableUserLoading = true;
            api.getTableUsers(this.table.id).then(res => {
                this.tableUserLoading = false;
                self.tableUsers = res.data;
            }, err => {
                this.tableUserLoading = false;
                this.$Message.error(err.response.data.message);
            })
        },
        updateColumnOrders() {
            var ids = this.columns.map(item => item.id);
            api.updateColumnOrders(this.table.id, {"ids":ids}).then(res => {
                this.$Message.success('修改列显示顺序成功!');
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        toggleFilters() {
            // 打开前重新设置下数据
            if(!this.showFilterDialog) {
                // debugger;
                this.resetFilters();
            }

            this.showFilterDialog = !this.showFilterDialog;

        },
        setFilterTypes(args) {
            // debugger;
            let self = this;
            this.filters.forEach((item, index) => {
                let cols = self.columns.filter(t => t.id == item.column_id);
                if(!cols || !cols.length) return;
                self.filters[index].conditions = self.conditions[cols[0].type];
                self.filters[index].conditions.forEach(t => {
                    // console.log(t.type, item.operator);
                    if(t.type == item.operator) {
                        self.filters[index].condition = t;
                    }
                });

            });
        },
        addFilter() {
            // debugger;
            var obj = {column_id:this.columns[0].id, value:""};
            obj.conditions = this.conditions[this.columns[0].type];
            obj.condition = obj.conditions[0];
            obj.operator = obj.condition.type;
            this.filters.push(obj);
        },
        removeFilter(index) {
            this.filters.splice(index,1);
        },
        resetFilters(list) {
            if(!list) list = this.view_filters;
            let self = this;
            let formated = list.map(item => {
                let obj = item;
                let cols = self.columns.filter(t => t.id == item.column_id);  
                if(cols.length) {
                    let conditions = self.conditions[cols[0].type] || [];
                    if(conditions.length) {
                        obj.conditions = conditions;
                        obj.conditions.forEach(t => {
                            if(t.type == item.operator) {
                                obj.condition = t;
                            }
                        });
                    }
                } 
                

                return obj;
            });
            this.filters = formated;
        },  
        updateFilters() {
            var _filters = [];
            this.filters.forEach(item => {
                _filters.push({column_id: item.column_id, operator: item.operator, value: item.value});
            });
            let _sorts = this.sorts.map(item => item.is_desc ? '-'+item.id : item.id);
            let self = this;
            let data = {'filters':_filters, 'filter_type': this.filterType};
            api.updateTableView(this.table.id, this.view.id, data)
            .then(res => {
                self.$store.dispatch('SetView', res.data);
                self.$emit('refreshData');
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },   
        toggleSorts() {
            // 打开前重新设置下数据
            if(!this.showSortsDialog) {

            }
            this.showSortsDialog = !this.showSortsDialog;
        },        
        updateSorts() {
            let _sorts = this.sorts.map(item => item.is_desc ? '-'+item.id : item.id);
            let self = this;
            let data = {'sorts':_sorts};
            api.updateTableView(this.table.id, this.view.id, data)
            .then(res => {
                self.$store.dispatch('SetView', res.data);
                self.$emit('refreshData');
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        resetSorts(list) {
            if(!list) list = this.view_sorts;
            let formated = list.map(item => {
                let obj = {};
                obj.is_desc = item[0] == '-';
                obj.id = obj.is_desc ? item.substring(1,item.length) : item;
                return obj;
            });
            this.sorts = formated;
        },
        addSort() {
            this.sorts.push({id:this.columns[0].id, is_desc: false});
            // this.formatedSorts.push({id:this.columns[0].id, is_desc: false});
        },
        removeSort(index) {
            this.sorts.splice(index, 1);
        },
        updateHiddenColumns() {
            let self = this;
            let data = {'hide_columns':this.hide_columns};
            api.updateTableView(this.table.id, this.view.id, data)
            .then(res => {
                self.$store.dispatch('SetView', res.data);
                self.$emit('refreshData');
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        openDeleteDialog() {
            if(!this.selectedRows.length) return;
            let ids = this.selectedRows.map(row => row.id);
            let self = this;
            this.$Modal.confirm({
                title: '确认',
                content: '<p>删除数据后无法再恢复, 确认要继续删除吗?</p>',
                width: 320,
                onOk: () => {
                    api.deleteRecords(self.table.id, ids).then(response => {
                        // let newRecords = self.records.filter(row => ids.indexOf(row.id) == -1);
                        // self.$store.dispatch('SetRecords', newRecords);
                        self.$emit('refreshData');
                        self.$Message.success("删除列成功");
                    }, err => {
                        this.$Message.error(err.response.data.message);
                        self.$emit('refreshData');
                    });        
                },
                onCancel: () => {
                    
                }
            });
        },
        showNotImplementWarning() {
            this.$Message.warning("此功能还未做好 :(");
        },
        onInput: debounce(function() {
            this.$store.dispatch('SetTableSearchText', this.searchText);
        }, 200),
        keyup(e) {
            // if(e.ctrlKey && e.keyCode==) {

            // }
        },
        refresh() {
            this.$store.commit(types.RECORDS, []);
            this.$emit('refreshData');
        }
    },
    watch: {
        view_sorts(newVal, oldVal) {
            this.resetSorts(newVal);
        },
        view_filters(newVal, oldVal) {
            // if(!this.conditions || !this.conditions.length) return;
            // this.resetFilters(newVal);
        }
    },
    created: function () {
        // 当Ctrl+F时, 自动对焦到搜索框里
        this.$store.dispatch('SetTableSearchText', "");
        window.addEventListener('keyup', this.keyup);
        this.getTableUsers();
    },
}
</script>
<style scoped>
.search-box {
    margin-top:8px;
    width: 150px;
    vertical-align:middle;
}
.toolbar {
    height: 55px;
    border-bottom: 1px solid #d7dde4;
    background: #f9faf9;
    /* padding: 3px 10px; */
    position: relative;
}
.action {
    float: left;
    text-align: center;
    cursor: pointer;
    padding: 4px 0;
    margin: 0 10px;
    border: 1px solid transparent;
}
.action[disabled] {
    color: #ccc;
    cursor: not-allowed;
}
.action.item {
    width: 55px;
}
.action.item:hover {
    background: #fff;
    border-color: #d8dcd8;
    box-shadow: 0 0 0.25rem hsla(0,0%,64%,.26);
}
.action .title {
    margin-top: 2px;
}
.right {
    float: right;
}
.toolbar hr {
    margin: 10px;
    height: 25px;
    width: 1px;
    background: #dde3dd;
    border: 0;
    box-sizing: content-box;
    float: left;
}
.dropdown-item {
    text-align: left;
    height: 30px;
    line-height: 30px;
}
.hide-columns .ivu-checkbox-wrapper {
    width: 120px;
}
.row-item {
     margin: 5px 0;
}
.sorts .ivu-select {
    width: 140px;
}

.sorts .row-item:hover {
    
}
.add-button {
    margin-top: 10px;
    color: #555;
    display:inline-block;
}
.draggable-item {
    height: 35px;
    line-height: 35px;
    cursor: pointer;
    padding: 0 10px;
    border-bottom: 1px solid #dedede;
}
.draggable-item:hover {
    background: #eee;
}
.table {
    border-spacing: 0;
    width:100%;
    border: 1px solid #d7dde4;
}
.table th {
    height: 40px;
    padding-left: 15px;
    white-space: nowrap;
    overflow:hidden;
    background: #f5f7f9;
    box-sizing: border-box;
    text-align: left;
    text-overflow: ellipsis;
    vertical-align: middle;
    border-bottom: 1px solid #e3e8ee;    
}
.table td {
    height: 40px;
    padding: 0 15px;
}
.user-table {
    width:100%;
}
.title-box {
    position: absolute;
    left: calc(50% - 100px);
    top: 0;
    color: #ccc;
    font-size: 14px;
    height: 100%;
    line-height: 55px;
    width: 200px;
    text-align: center;
}


</style>
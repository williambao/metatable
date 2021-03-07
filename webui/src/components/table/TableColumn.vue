<template>
<div class="column" :style="{width:column.width+'px', 'text-align': column.text_align}">
    <div class="cell" v-for="record in records" :key="'cell_'+column.id+'_'+record.id">
        <span>
            {{record.cells[column.id]}}
        </span>
    </div>
</div>
</template>
<script>
import api from '../../axios'
export default {
    props: ["column", "records", "activeRow", "isActive"],
}
</script>
<style>
.column {
    float: left;
    border-right: 1px solid #d7dde4;
}
.cell {
    height: 40px;
    border-bottom: 1px solid #d7dde4;
}
</style>
<template>
  <div>
    <t-pagination
      v-model="current"
      :total="36"
      :page-size.sync="pageSize"
      @current-change="onCurrentChange"
      @page-size-change="onPageSizeChange"
      @change="onChange"
    />
    <br /><br />
    
    <t-pagination
      showFirstAndLastPageBtn
      :total="36"
      :default-current="2"
      :default-page-size="10"
      @current-change="onCurrentChange"
      @page-size-change="onPageSizeChange"
      @change="onChange"
    />
  </div>
</template>

<script>
export default {
  data() {
    return {
      current: 12,
      pageSize: 5,
    };
  },
  methods: {
    onPageSizeChange(size, pageInfo) {

    },
    onCurrentChange(current, pageInfo) {

    },
    onChange(pageInfo) {

    },
  },
};
</script>
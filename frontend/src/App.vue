<template>
  <el-container>
    <el-header>
      <el-row type="flex" class="row-bg" justify="space-between">
        <el-col :span="6"><el-image :src="src" class="logo" /></el-col>
        <el-col :span="6">
          <el-button type="primary" icon="el-icon-plus" circle @click="addRow()"></el-button>
          <el-button type="success" icon="el-icon-check" circle @click="apply()"></el-button>
        </el-col>
      </el-row>
    </el-header>
    <el-divider></el-divider>
    <el-main>
      <el-table
        :data="tableData"
        :row-class-name="tableRowClass"
        border
        style="width: 100%"
      >
        <el-table-column label="IP">
          <template v-slot="{ row }">
            <template v-if="row.edit">
              <el-input v-model="row.ip" class="edit-input" size="small" />
            </template>
            <span v-else>{{ row.ip }}</span>
          </template>
        </el-table-column>
        <el-table-column label="FQDN">
          <template v-slot="{ row }">
            <template v-if="row.edit">
              <el-input v-model="row.fqdn" class="edit-input" size="small" />
            </template>
            <span v-else>{{ row.fqdn }}</span>
          </template>
        </el-table-column>
        <el-table-column label="Actions">
          <template v-slot="{ row }">
            <template v-if="row.edit">
              <el-button
                type="success"
                icon="el-icon-check"
                size="mini"
                circle
                @click="submitRow(row)"
              ></el-button>
              <el-button icon="el-icon-close" size="mini" circle @click="cancelEdit(row)"></el-button>
            </template>
            <el-button
              v-else
              type="primary"
              icon="el-icon-edit"
              size="mini"
              circle
              @click="startEdit(row)"
            ></el-button>
            <el-popconfirm
              confirmButtonText="OK"
              cancelButtonText="Cancel"
              icon="el-icon-info"
              iconColor="red"
              title="Are you sure you want to delete this item?"
              @confirm="confirmDelete(row)"
            >
              <template #reference>
                <el-button
                  type="danger"
                  icon="el-icon-delete"
                  size="mini"
                  circle
                ></el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-main>
  </el-container>
</template>
<script>
import { ref } from "vue";
import logo from "./assets/logo.png";

export default {
  name: "App",
  components: {},
  methods: {
    startEdit(row) {
      row.edit = true;
      console.log(row.ip);
    },
    cancelEdit(row) {
      row.edit = false;
      console.log(row.ip);
    },
    submitRow(row) {
      console.log(row);
      row.edit = false;
    },
    confirmDelete(row) {
      console.log("Delete", row);
    },
    addRow() {
      this.tableData.unshift({
        edit: true
      });
    },
    async apply() {

    },
    tableRowClass({ row }) {
      if (row.edit) {
        return "edit-row";
      }
      return "";
    },
  },
  data() {
    return {
      src: logo,
      tableData: [
        {
          ip: "1.1.1.1",
          fqdn: "baidu.com",
          edit: false,
        },
      ],
    };
  },
  setup() {
    const msg = ref(
      "Welcome to Element Plus, a Vue 3.0 based component library"
    );
    const startHacking = () => {
      msg.value = "Start coding with Element Plus with 💖";
    };
    return {
      msg,
      startHacking,
    };
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
.logo {
  width: 50%;
}
.edit-row > td,
.edit-row:hover > td {
  background-color: #fdf6ec !important;
}
.btn {
  margin-top: 100px;
}
</style>

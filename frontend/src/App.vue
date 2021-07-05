<template>
  <el-container>
    <el-header>
      <el-row type="flex" class="row-bg" justify="space-between">
        <el-col :span="6"><el-image :src="src" class="logo" /></el-col>
        <el-col :span="6">
          <el-button
            type="primary"
            icon="el-icon-plus"
            circle
            :disabled="currentRow != undefined"
            @click="addRow()"
          ></el-button>
          <el-button
            type="success"
            icon="el-icon-check"
            circle
            @click="apply()"
          ></el-button>
        </el-col>
      </el-row>
    </el-header>
    <el-divider></el-divider>
    <el-main>
      <el-table
        :data="tableData"
        :row-class-name="tableRowClass"
        v-loading="loading"
        stripe="true"
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
        <el-table-column label="Comment">
          <template v-slot="{ row }">
            <template v-if="row.edit">
              <el-input v-model="row.comment" class="edit-input" size="small" />
            </template>
            <span v-else>{{ row.comment }}</span>
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
                native-type="submit"
                @click="submitRow(row)"
              ></el-button>
              <el-button
                icon="el-icon-close"
                size="mini"
                circle
                @click="cancelEdit()"
              ></el-button>
            </template>
            <template v-if="currentRow == undefined">
              <el-button
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
          </template>
        </el-table-column>
      </el-table>
    </el-main>
  </el-container>
</template>
<script>
import logo from "./assets/logo.png";
import hosts from "./apis/hosts";

export default {
  name: "App",
  components: {},
  methods: {
    handleError(error) {
      console.log(error)
      this.$notify({
        title: "Error",
        message: error.message,
        type: "error",
        offset: 50,
      });
    },
    startEdit(row) {
      row.edit = true;
      this.currentRow = Object.assign({}, row);
    },
    cancelEdit() {
      if (this.currentRow.id == undefined) {
        this.tableData.shift();
      }
      else {
        this.currentRow.edit = false;
        this.tableData[this.tableData.findIndex(item => item.id == this.currentRow.id)] = this.currentRow;
      }
      this.currentRow = undefined;
    },
    async submitRow(row) {
      // undefined id means this is a new row
      if (this.currentRow.id == undefined) {
        hosts.createHost(row).then(response => {
          this.tableData = this.tableData.filter(item => item.id != undefined);
          this.tableData.unshift(response.data);
          this.$notify({
            title: "Success",
            message: "Create successfully",
            type: "success",
            offset: 50,
          });
        }).catch(error => this.handleError(error))
      }
      else {
        hosts.updateHost(row).then(response => {
          row = response.data;
          this.tableData[this.tableData.findIndex(item => item.id == row.id)] = row;
          this.$notify({
            title: "Success",
            message: "Update successfully",
            type: "success",
            offset: 50,
          });
        }).catch(error => this.handleError(error))
      }
      row.edit = false;
      this.currentRow = undefined;
    },
    async confirmDelete(row) {
      hosts.deleteHost(row.id).then(
        () => {
          this.tableData = this.tableData.filter(item => item.id != row.id);
          this.$notify({
            title: "Success",
            message: "Delete successfully",
            type: "success",
            offset: 50,
          });
        }
      ).catch(error => this.handleError(error));
    },
    addRow() {
      this.currentRow = {
        edit: true,
      };
      this.tableData.unshift(this.currentRow);
    },
    async apply() {
      hosts.applyChange().then(() => {
        this.$notify({
          title: "Success",
          message: "Apply DNS successfully",
          type: "success",
          offset: 50,
        });
      }).catch(error => this.handleError(error));
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
      tableData: [],
      loading: false,
      currentRow: undefined,
    };
  },
  created() {
    this.loading = true;
    hosts
      .getHosts()
      .then((response) => (this.tableData = response.data))
      .catch(response => this.handleError(response))
      .finally(() => {
        this.loading = false;
      });
  },
  setup() {
    return {};
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
.el-main {
  padding: 0;
}
.row-bg {
  height: 50px;
}
.logo {
  width: 50%;
}

.logo > img {
  height: 80px;
  width: auto;
}

.edit-row > td,
.edit-row:hover > td {
  background-color: #fdf6ec !important;
}
.btn {
  margin-top: 100px;
}
</style>

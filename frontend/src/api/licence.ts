import store from "@/store/store";
import axios from "axios";

export const getLicenses = async () => {
  try {
    const res = await axios.get("http://127.0.0.1:8080/api/admin/licenses", {
      headers: {
        Authorization: "JWT " + localStorage.getItem("token"),
      },
    });
    store.changeLicenses(res.data)
    console.log(res.data);
    
  } catch (error) {
    console.log(error);
    return [];
  }
};

export const createLicense = async () => {
  try {
    const now = new Date();
    const nextMonth = new Date(now);
    nextMonth.setMonth(now.getMonth() + 1);

    await axios.post('http://127.0.0.1:8080/api/license/generate', {
      "access_level": {
    "max_messages": +store.maxMessages,
    "max_users": +store.maxUsers
  },
  "expiration_date": nextMonth,
  "owner": store.owner
    }, {
      headers: {
        Authorization: "JWT " + localStorage.getItem("token"),
      },
    })
    getLicenses()
  } catch (error) {
    console.log(error);
    
  }
}

export const downloadLicense = async (id: string) => {
  try {
    const response = await axios.get(`http://127.0.0.1:8080/api/admin/licenses/${id}/download`, {
      headers: {
        Authorization: "JWT " + localStorage.getItem("token"),
      },
      responseType: "blob", // 👈 важно!
    });

    // Создание blob и временной ссылки
    const blob = new Blob([response.data], { type: "application/octet-stream" });
    const url = window.URL.createObjectURL(blob);

    const link = document.createElement("a");
    link.href = url;
    link.download = `license-${id}.lis`; // 👈 имя файла

    document.body.appendChild(link);
    link.click();

    // Очистка
    link.remove();
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error("Ошибка при загрузке лицензии:", error);
  }
};


export const deleteLicense = async (id: string) => {
  try {
     await axios.delete(`http://127.0.0.1:8080/api/admin/licenses/${id}`, {
      headers: {
        Authorization: 'JWT ' + localStorage.getItem('token')
      }
     })
     getLicenses()
  } catch (error) {
    console.log(error);
  }
}
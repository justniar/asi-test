import axios from "axios";

const API_URL = 'http://localhost:8080/api/clients';

export const createClient = async (formData: FormData) => {
    const response = await axios.post(API_URL, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
};

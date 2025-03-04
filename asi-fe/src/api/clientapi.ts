import axios from 'axios';

const API_URL = 'http://localhost:8080/api/clients';

export const getClients = async () => {
    const response = await axios.get(API_URL);
    return response.data;
};

export const getClient = async (slug: string) => {
    const response = await axios.get(`${API_URL}/${slug}`);
    return response.data;
};

export const createClient = async (data: any) => {
    const response = await axios.post(API_URL, data);
    return response.data;
};

export const updateClient = async (slug: string, data: any) => {
    const response = await axios.put(`${API_URL}/${slug}`, data);
    return response.data;
};

export const deleteClient = async (slug: string) => {
    await axios.delete(`${API_URL}/${slug}`);
};

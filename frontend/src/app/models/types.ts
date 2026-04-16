export interface Product {
    id?: number;
    code: string;
    description: string;
    balance: number;
}

export interface InvoiceItem {
    id?: number;
    product_code: string;
    quantity: number;
}

export interface Invoice {
    id?: number;
    number?: number;
    status: 'Aberta' | 'Fechada';
    items: InvoiceItem[];
    created_at?: string;
}

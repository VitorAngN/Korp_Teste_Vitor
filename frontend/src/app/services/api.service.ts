import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError, BehaviorSubject } from 'rxjs';
import { catchError, map, tap } from 'rxjs/operators';
import { Product, Invoice } from '../models/types';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private stockUrl = 'http://localhost:8081/api';
  private invoiceUrl = 'http://localhost:8082/api';

  private productsSubject = new BehaviorSubject<Product[]>([]);
  public products$ = this.productsSubject.asObservable();
  
  private invoicesSubject = new BehaviorSubject<Invoice[]>([]);
  public invoices$ = this.invoicesSubject.asObservable();

  constructor(private http: HttpClient) { }

  loadProducts(): Observable<Product[]> {
    return this.http.get<Product[]>(`${this.stockUrl}/products`).pipe(
      tap(products => this.productsSubject.next(products)),
      catchError(this.handleError)
    );
  }

  createProduct(product: Product): Observable<Product> {
    return this.http.post<Product>(`${this.stockUrl}/products`, product).pipe(
      tap(() => this.loadProducts().subscribe()), 
      catchError(this.handleError)
    );
  }

  generateDescriptionAI(productName: string): Observable<any> {
    return this.http.post(`${this.stockUrl}/products/ai/generate`, { product_name: productName }).pipe(
      catchError(this.handleError)
    );
  }

  loadInvoices(): Observable<Invoice[]> {
    return this.http.get<Invoice[]>(`${this.invoiceUrl}/invoices`).pipe(
      tap(invoices => this.invoicesSubject.next(invoices)),
      catchError(this.handleError)
    );
  }

  createInvoice(invoice: Partial<Invoice>): Observable<Invoice> {
    return this.http.post<Invoice>(`${this.invoiceUrl}/invoices`, invoice).pipe(
      tap(() => this.loadInvoices().subscribe()),
      catchError(this.handleError)
    );
  }

  printInvoice(id: number, simulateFailure: boolean = false): Observable<any> {
    const headers = { 'X-Simulate-Failure': simulateFailure ? 'true' : 'false' };
    return this.http.post(`${this.invoiceUrl}/invoices/${id}/print`, {}, { headers }).pipe(
      tap(() => {
        this.loadInvoices().subscribe(); 
        this.loadProducts().subscribe(); 
      }),
      catchError(this.handleError)
    );
  }

  private handleError(error: HttpErrorResponse) {
    let errorMessage = 'Ocorreu um erro desconhecido.';
    // Erros do Microsserviço vêm via json no body
    if (error.error && error.error.error) {
      errorMessage = error.error.error;
    } else if (error.message) {
      errorMessage = error.message;
    }
    return throwError(() => new Error(errorMessage));
  }
}

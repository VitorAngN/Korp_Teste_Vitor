import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormGroup, FormArray, FormBuilder, Validators, ReactiveFormsModule } from '@angular/forms';
import { ApiService } from '../../services/api.service';
import { Invoice, Product } from '../../models/types';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-invoices',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './invoices.component.html',
  styleUrl: './invoices.component.css'
})
export class InvoicesComponent implements OnInit, OnDestroy {
  invoices: Invoice[] = [];
  products: Product[] = [];
  
  invoiceForm: FormGroup;
  isSubmitting = false;
  printProcessing: { [key: number]: boolean } = {};
  toastMsg = '';
  isErrorToast = false;
  
  private destroy$ = new Subject<void>();

  constructor(private fb: FormBuilder, private api: ApiService) {
    this.invoiceForm = this.fb.group({
      items: this.fb.array([])
    });
  }

  ngOnInit() {
    this.api.invoices$
      .pipe(takeUntil(this.destroy$))
      .subscribe(data => this.invoices = data);

    this.api.products$
      .pipe(takeUntil(this.destroy$))
      .subscribe(data => this.products = data);
      
    this.api.loadInvoices().subscribe();
    this.api.loadProducts().subscribe();
    this.addItem();
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  get items() {
    return this.invoiceForm.get('items') as FormArray;
  }

  addItem() {
    this.items.push(this.fb.group({
      product_code: ['', Validators.required],
      quantity: [1, [Validators.required, Validators.min(1)]]
    }));
  }

  removeItem(index: number) {
    this.items.removeAt(index);
  }

  onSubmit() {
    if (this.invoiceForm.invalid || this.items.length === 0) return;
    
    this.isSubmitting = true;
    const values = this.invoiceForm.value;

    this.api.createInvoice({ items: values.items }).subscribe({
      next: () => {
        this.items.clear();
        this.addItem();
        this.isSubmitting = false;
        this.showToast('Nota Fiscal Aberta com Sucesso!');
      },
      error: (err) => {
        this.isSubmitting = false;
        this.showToast(err.message, true);
      }
    });
  }

  printInvoice(invoice: Invoice) {
    if (invoice.status !== 'Aberta' || !invoice.id) {
       this.showToast("Só é possível imprimir notas com status 'Aberta'.", true);
       return;
    }

    this.printProcessing[invoice.id] = true;

    this.api.printInvoice(invoice.id).subscribe({
      next: () => {
        this.printProcessing[invoice.id!] = false;
        this.showToast(`Nota Fiscal #${invoice.number || invoice.id} fechada e impressa! Os saldos foram abatidos.`);
      },
      error: (err) => {
        this.printProcessing[invoice.id!] = false;
        this.showToast(`Erro na Impressão: ${err.message}`, true);
      }
    });
  }

  showToast(msg: string, error: boolean = false) {
    this.toastMsg = msg;
    this.isErrorToast = error;
    setTimeout(() => this.toastMsg = '', 6000);
  }
}

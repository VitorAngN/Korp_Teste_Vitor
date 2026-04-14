import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormGroup, FormControl, Validators, ReactiveFormsModule } from '@angular/forms';
import { ApiService } from '../../services/api.service';
import { Product } from '../../models/types';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-products',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './products.component.html',
  styleUrl: './products.component.css'
})
export class ProductsComponent implements OnInit, OnDestroy {
  products: Product[] = [];
  productForm = new FormGroup({
    code: new FormControl('', Validators.required),
    description: new FormControl('', Validators.required),
    balance: new FormControl(0, [Validators.required, Validators.min(0)])
  });
  
  isSubmitting = false;
  toastMsg = '';
  
  private destroy$ = new Subject<void>();

  constructor(private api: ApiService) {}

  ngOnInit() {
    this.api.products$
      .pipe(takeUntil(this.destroy$))
      .subscribe(data => this.products = data);
    
    this.api.loadProducts().subscribe();
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  onSubmit() {
    if (this.productForm.invalid) return;
    
    this.isSubmitting = true;
    const prod: Product = {
      code: this.productForm.value.code!,
      description: this.productForm.value.description!,
      balance: this.productForm.value.balance!
    };

    this.api.createProduct(prod).subscribe({
      next: () => {
        this.productForm.reset({ balance: 0 });
        this.isSubmitting = false;
        this.showToast('Produto cadastrado com sucesso!');
      },
      error: (err) => {
        console.error(err);
        this.isSubmitting = false;
        this.showToast(err.message, true);
      }
    });
  }

  isErrorToast = false;
  
  showToast(msg: string, error: boolean = false) {
    this.toastMsg = msg;
    this.isErrorToast = error;
    setTimeout(() => this.toastMsg = '', 4000);
  }
}

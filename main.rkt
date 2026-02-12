#lang racket


(define rmember
  (lambda (item l)
    (cond
      [(null? l) '()]
      [(equal? item (car l)) (rmember item (cdr l))]
      [else (cons (car l) (rmember item (cdr l)))])))

(print (rmember 'a '(b c d a e)))

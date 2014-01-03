(define not
  (lambda (b)
    (if b
        #f
        #t)))

(define and
  (lambda l
    (if (null? l)
        #t
        (if (not (car l))
            #f
            (and (cdr l))))))
            

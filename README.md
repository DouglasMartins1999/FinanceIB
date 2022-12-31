# FinanceIB

#### Gerador Automático de Demonstrativo Financeiro da IBR SJ

Aplicativo escrito em [Go](https://pt.wikipedia.org/wiki/Go_(linguagem_de_programa%C3%A7%C3%A3o)) utilizando as bibliotecas [ofxgo](https://github.com/aclindsa/ofxgo), [excelize](https://github.com/qax-os/excelize) e [maroto](https://github.com/johnfercher/maroto).

Recebe um extrato bancário no formato OFX e gera uma planilha XLSX a qual o usuário pode editar a descrição, data e valor de cada entrada do extrato, bem como remover ou adicionar uma entrada. A planilha também permite customizar mês, ano e data de assinatura do relatório. Ao passar a planilha novamente ao aplicativo, será gerado um demonstrativo no formato PDF.
package controllers

import (
	"GoSecKill/internal/services"
	"GoSecKill/pkg/models"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"go.uber.org/zap"
)

type ProductController struct {
	productService services.IProductService
	orderService   services.IOrderService
	sessions       *sessions.Sessions
}

func NewProductController(productService services.IProductService, orderService services.IOrderService, sessions *sessions.Sessions) *ProductController {
	return &ProductController{productService: productService, orderService: orderService, sessions: sessions}
}

var (
	htmlOutPath = "./fronted/web/htmlProductShow/"

	templatePath = "./fronted/web/views/template/"
)

func (p *ProductController) GetGenerateHtml(ctx iris.Context) {
	productString := ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		zap.L().Error("Failed to convert productID", zap.Error(err))
	}

	contentTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		zap.L().Error("Failed to parse template", zap.Error(err))
	}

	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")

	product, err := p.productService.GetProductByID(productID)
	if err != nil {
		zap.L().Error("Failed to get product", zap.Error(err))
	}

	generateStaticHtml(ctx, contentTmp, fileName, &product)
}

func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *models.Product) {
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			zap.L().Error("Failed to remove file", zap.Error(err))
		}
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		zap.L().Error("Failed to open file", zap.Error(err))
	}
	defer file.Close()
	template.Execute(file, &product)
}

func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func (p *ProductController) GetDetail(ctx iris.Context) mvc.View {
	product, err := p.productService.GetProductByID(1)
	if err != nil {
		zap.L().Error("Failed to get product", zap.Error(err))
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder(ctx iris.Context) mvc.View {
	productString := ctx.URLParam("productID")
	userString := ctx.GetCookie("uid")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		zap.L().Error("Failed to convert productID", zap.Error(err))
	}
	product, err := p.productService.GetProductByID(productID)
	if err != nil {
		zap.L().Error("Failed to get product", zap.Error(err))
	}
	var orderID uint
	showMessage := "Failed to buy product, please try again later"

	if product.Stock > 0 {
		product.Stock -= 1
		err := p.productService.UpdateProduct(product)
		if err != nil {
			zap.L().Error("Failed to update product", zap.Error(err))
		}
		//创建订单
		userID, err := strconv.Atoi(userString)
		if err != nil {
			zap.L().Error("Failed to convert userID", zap.Error(err))
		}

		order := &models.Order{
			UserId:     uint(userID),
			ProductId:  uint(productID),
			ProductNum: 1,
			Status:     models.OrderStatusSuccess,
		}

		order, err = p.orderService.InsertOrder(*order)
		orderID = order.ID
		if err != nil {
			zap.L().Error("Failed to insert order", zap.Error(err))
		} else {
			showMessage = "Successfully buy product"
		}
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": showMessage,
		},
	}

}

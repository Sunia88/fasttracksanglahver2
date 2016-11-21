$(document).ready(function(){
	var nocm = $("#nocm").val()
	
	$("#nocm").focus(function(){
		var value = $("#nocm").val();
		if (value == ""){
			$("#datapasien").html("Masukkan No. CM");
			nocm = "";
		} else {
			$("#datapasien").html("");
			nocm = "";
		}
		
	});
	
	$("#form1").on("keyup", "#nocm", function(){
		
		var value = $("#nocm").val();
		
		if (value == ""){
			$("#datapasien").html("Masukkan No. CM");
			nocm = "";
		} else if (value.length < 8){
			$("#datapasien").html("No. CM tidak lengkap");
			nocm = "";
		} else {
			$("#datapasien").html("Mantaf");
			nocm = "";
		}
	})
	
	
	
})